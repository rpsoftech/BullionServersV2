package auth

// func apiAdminLogin(c *fiber.Ctx) error {
// 	body := new(adminLoginBody)
// 	c.BodyParser(body)
// 	if errs := validator.Validator.Validate(body); len(errs) > 0 {
// 		err := &interfaces.RequestError{
// 			StatusCode: 400,
// 			Code:       interfaces.ERROR_INVALID_INPUT,
// 			Message:    "",
// 			Name:       "INVALID_INPUT",
// 			Extra:      errs,
// 		}
// 		err.AppendValidationErrors(errs)
// 		return err
// 	}
// 	entity, err := services.GeneralUserService.AdminLogin(body.Id, body.Password)
// 	if err != nil {
// 		return err
// 	} else {
// 		return c.JSON(entity)
// 	}
// }
