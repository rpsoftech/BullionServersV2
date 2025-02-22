package services

import (
	"net/http"

	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
)

type orderGeneralService struct {
	eventBus        *eventBusService
	firebaseDb      *firebaseDatabaseService
	bullionService  *bullionDetailsService
	flagService     *FlagServiceStruct
	liveRateService *liveRateServiceStruct
	groupService    *tradeUserGroupServiceStruct
	orderRepo       *repos.OrderRepoStruct
	bankRateService *bankRateService
	userService     *tradeUserServiceStruct
	productService  *productService
}

var OrderGeneralService *orderGeneralService

func init() {
	getOrderGeneralService()
}
func getOrderGeneralService() *orderGeneralService {
	if OrderGeneralService == nil {
		OrderGeneralService = &orderGeneralService{
			eventBus:        getEventBusService(),
			firebaseDb:      getFirebaseRealTimeDatabase(),
			bullionService:  getBullionService(),
			groupService:    getTradeUserGroupService(),
			flagService:     getFlagService(),
			liveRateService: getLiveRateService(),
			userService:     getTradeUserService(),
			bankRateService: getBankRateService(),
			productService:  getProductService(),
			orderRepo:       repos.OrderRepo,
		}
		println("Order General Service Initialized")
	}
	return OrderGeneralService
}

func (service *orderGeneralService) ValidateUserGroupForTrade(group *interfaces.TradeUserGroupEntity) (bool, error) {
	// Check for Group Activation
	if !group.IsActive {
		return false, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Group Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}
	// Check for Group User Can Trade
	if !group.CanTrade {
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED_FOR_GROUP,
			Message:    "Trading is disabled for your group. Contact User",
			Name:       "GROUP_NOT_ACTIVE",
		}
	}
	return true, nil
}

func (service *orderGeneralService) ValidateUserAndGroupMapWithWeight(groupMap *interfaces.TradeUserGroupMapEntity, weight int) (bool, error) {
	// Check for Group Map Activation
	if !groupMap.IsActive {
		return false, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Group Map Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}

	// Check for Group Map Can Trade
	if !groupMap.CanTrade {
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED_FOR_PRODUCT,
			Message:    "Trading is disabled for your group map. Contact Admin",
			Name:       "GROUP_NOT_ACTIVE",
		}
	}

	return service.validateVolumeForGroupMap(groupMap, weight)
}

func (service *orderGeneralService) validateVolumeForGroupMap(groupMap *interfaces.TradeUserGroupMapEntity, weight int) (bool, error) {
	if !groupMap.ValidateVolume(weight) {
		return false, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_VOLUME,
			Message:    "Invalid Volume",
			Name:       "INVALID_VOLUME",
		}
	}
	return true, nil
}

func (service *orderGeneralService) findOrderDetailsAndValidate(userId string, groupId string, groupMapId string, weight int) (*interfaces.TradeUserEntity, *interfaces.TradeUserGroupEntity, *interfaces.TradeUserGroupMapEntity, error) {
	// Get User
	user, err := service.userService.GetTradeUserById(userId)
	if err != nil {
		return nil, nil, nil, err
	}

	if user.GroupId != groupId {
		return nil, nil, nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "MissMatch Group Id",
			Name:       "MISS_MATCH_GROUP_ID",
			Extra:      "Solution Logout And Relogin",
		}
	}
	// Check for User Activation
	if !user.IsActive {
		return nil, nil, nil, &interfaces.RequestError{
			StatusCode: http.StatusUnauthorized,
			Code:       interfaces.ERROR_PERMISSION_NOT_ALLOWED,
			Message:    "Account Is Not Active Please Contact Admin",
			Name:       "ERROR_PERMISSION_NOT_ALLOWED",
		}
	}

	if flags, err := service.flagService.GetFlags(user.BullionId); err != nil {
		return nil, nil, nil, err
	} else if !flags.CanTrade {
		// Check If Trading Is Disabled
		return nil, nil, nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_TRADING_IS_DISABLED,
			Message:    "Trading is disabled. Contact User",
			Name:       "BULLION_NOT_ACTIVE",
		}
	}
	// Get Group
	group, err := service.groupService.GetGroupByGroupId(groupId, user.BullionId)
	if err != nil {
		return nil, nil, nil, err
	}
	if valid, err := service.ValidateUserGroupForTrade(group); !valid && err != nil {
		return nil, nil, nil, err
	}
	// Get Group Map
	groupMaps, err := service.groupService.GetGroupMapByGroupId(groupId, user.BullionId)
	if err != nil {
		return nil, nil, nil, err
	}
	var groupMap *interfaces.TradeUserGroupMapEntity

	for _, v := range *groupMaps {
		if v.ID == groupMapId {
			groupMap = &v
			break
		}
	}
	if groupMap == nil {
		return nil, nil, nil, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GROUP_MAP_NOT_FOUND,
			Message:    "Group Map Not Found",
			Name:       "GROUP_MAP_NOT_FOUND",
			Extra:      "Solution Logout And Relogin",
		}

	}
	if valid, err := service.ValidateUserAndGroupMapWithWeight(groupMap, weight); !valid && err != nil {
		return nil, nil, nil, err
	}

	return user, group, groupMap, nil
}
func (service *orderGeneralService) PlaceOrder(orderType interfaces.OrderStatus, userId string, groupId string, groupMapId string, buySell interfaces.BuySell, weight int, price float64, placedBy string) (*interfaces.OrderEntity, error) {

	user, group, groupMap, err := service.findOrderDetailsAndValidate(userId, groupId, groupMapId, weight)
	if err != nil {
		return nil, err
	}
	// product
	product, err := service.productService.GetProductsById(group.BullionId, groupMap.ProductId)
	if err != nil {
		return nil, err
	}

	// Validate Pricing
	finalRate, err := service.calCulateAndReturnFinalRateForOrder(product, group, groupMap, buySell)
	println("Final Rate", finalRate)
	order := &interfaces.OrderEntity{
		BaseEntity: &interfaces.BaseEntity{},
		OrderBase: &interfaces.OrderBase{
			BullionId:   group.BullionId,
			OrderType:   interfaces.OrderType(orderType),
			BuySell:     buySell,
			ProductName: product.Name,
		},
		LimitWatcherRequired: &interfaces.LimitWatcherRequired{
			ProductId:         product.ID,
			GroupId:           group.ID,
			ProductGroupMapId: groupMap.ID,
			Volume:            float64(weight),
			Weight:            weight,
		},
		// DeliveryData: ,
		// Identity: ,
		// AfterSuccessOrder: ,
	}
	println("Order", order)
	// order.BaseEntity.
	// TODO Check Hedging And Place Order

	// TODO Update Order Entity in DB

	_, err = user.UpdateMarginAfterOrder(weight, product.SourceSymbol)
	if err != nil {
		return nil, err
	}

	// TODO Check Hedging And Place Order
	return nil, nil
	// return service.orderRepo.PlaceOrder(orderType, user, group, groupMap, price, placedBy)
	// return service.orderRepo.PlaceOrder(orderType, user, group, groupMap, price, placedBy)
}

func (service *orderGeneralService) calCulateAndReturnFinalRateForOrder(product *interfaces.ProductEntity, group *interfaces.TradeUserGroupEntity, groupMap *interfaces.TradeUserGroupMapEntity, buySell interfaces.BuySell) (float64, error) {
	priceReadKey := interfaces.PRICE_ASK
	if product.CalculatedOnPriceOf == interfaces.CALCULATE_ON_BID {
		priceReadKey = interfaces.PRICE_BID
	} else if product.CalculatedOnPriceOf == interfaces.CALCULATE_ON_BID_ASK {
		if buySell == interfaces.Sell {
			priceReadKey = interfaces.PRICE_ASK
		} else {
			priceReadKey = interfaces.PRICE_BID
		}
	}

	productSymbol := product.SourceSymbol.ToSymbolEnum()

	groupPremium := group.Gold

	if product.CalcPriceMethod == interfaces.CALCULATION_PRICE_TYPE_BANK {
		if product.SourceSymbol == interfaces.SOURCE_SYMBOL_GOLD {
			productSymbol = interfaces.SYMBOL_GOLD_SPOT
		} else {
			productSymbol = interfaces.SYMBOL_SILVER_SPOT
			groupPremium = group.Silver
		}
	}

	rate := service.liveRateService.GetLiveRate(productSymbol, priceReadKey)

	if rate == 0 {
		return 0, &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_LIVE_RATE_NOT_FOUND,
			Message:    "Live Rate Not Found",
			Name:       "LIVE_RATE_NOT_FOUND",
		}
	}

	if product.CalcPriceMethod == interfaces.CALCULATION_PRICE_TYPE_BANK {
		bankRate, err := service.bankRateService.GetBankRateCalcByBullionId(group.BullionId)
		if err != nil {
			return 0, err
		}
		inrRate := service.liveRateService.GetLiveRate(interfaces.SYMBOL_INR, priceReadKey)
		calcFunc := bankRate.GOLD_SPOT.CalculatePrice
		if product.SourceSymbol == interfaces.SOURCE_SYMBOL_SILVER {
			calcFunc = bankRate.SILVER_SPOT.CalculatePrice
		}
		rate = calcFunc(rate, inrRate)
		// rate = bankRate.Rate
	}

	// Extra Premium For Group
	extraPremium := groupMap.Sell + groupPremium.Sell

	calcSnapshot := &product.CalcSnapshot.Sell

	if buySell == interfaces.Buy {
		calcSnapshot = &product.CalcSnapshot.Buy
		extraPremium = groupMap.Buy + groupPremium.Buy
	}

	finalRate := interfaces.Calculate(rate+extraPremium, calcSnapshot)
	return finalRate, nil
	// return service.orderRepo.GetOrderById(orderId)
}

// check user is valid
// check group is valid
// check group map is valid
// check for volume
