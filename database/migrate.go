package database

import (
	"log"

	"loan-api/models"
)

// Migrate ejecuta las migraciones de la base de datos
func Migrate() error {
	// Ejecutar migraciones automáticas
	err := DB.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.LoanType{},
		&models.LoanTypeVersion{},
		&models.LoanTypeForm{},
		&models.LoanTypeVersionFormInput{},
		&models.Loan{},
		&models.LoanData{},
	)

	if err != nil {
		log.Printf("Error en la migración: %v", err)
		return err
	}

	// Ejecutar seeders para datos iniciales
	if err := seedInitialData(); err != nil {
		log.Printf("Error en los seeders: %v", err)
		return err
	}

	log.Println("Migraciones ejecutadas correctamente")
	return nil
}

// seedInitialData inserta datos iniciales de configuración
func seedInitialData() error {
	// Crear tenant de prueba
	var tenant models.Tenant
	result := DB.Where("code = ?", "test_bank").First(&tenant)
	if result.Error != nil {
		// El tenant no existe, crearlo
		tenant = models.Tenant{
			Name:        "Banco de Pruebas",
			Code:        "test_bank",
			Description: "Entidad crediticia para pruebas de desarrollo",
			IsActive:    true,
			Config:      `{"max_loan_amount": 50000000, "min_credit_score": 500}`,
		}
		if err := DB.Create(&tenant).Error; err != nil {
			return err
		}
	}

	// Crear tipo de préstamo de prueba
	var loanType models.LoanType
	result = DB.Where("tenant_id = ? AND code = ?", tenant.ID, "personal_loan").First(&loanType)
	if result.Error != nil {
		loanType = models.LoanType{
			TenantID:    tenant.ID,
			Name:        "Préstamo Personal",
			Code:        "personal_loan",
			Description: "Préstamo personal para gastos diversos",
			IsActive:    true,
			MinAmount:   100000,
			MaxAmount:   10000000,
		}
		if err := DB.Create(&loanType).Error; err != nil {
			return err
		}
	}

	// Crear versión del tipo de préstamo
	var loanTypeVersion models.LoanTypeVersion
	result = DB.Where("loan_type_id = ? AND version = ?", loanType.ID, "1.0").First(&loanTypeVersion)
	if result.Error != nil {
		loanTypeVersion = models.LoanTypeVersion{
			LoanTypeID:  loanType.ID,
			Version:     "1.0",
			Description: "Versión inicial del préstamo personal",
			IsActive:    true,
			IsDefault:   true,
			Config:      `{"approval_rules": {"min_income": 1000000, "max_debt_ratio": 0.4}}`,
		}
		if err := DB.Create(&loanTypeVersion).Error; err != nil {
			return err
		}
	}

	// Verificar si ya existen formularios para esta versión
	var existingFormsCount int64
	DB.Model(&models.LoanTypeForm{}).Where("loan_type_version_id = ?", loanTypeVersion.ID).Count(&existingFormsCount)

	if existingFormsCount == 0 {
		// Crear todos los formularios en una sola operación
		forms := []models.LoanTypeForm{
			{
				LoanTypeVersionID: loanTypeVersion.ID,
				Label:             "Información Personal",
				Code:              "personal_info",
				Description:       "Datos personales del solicitante",
				Order:             1,
				IsRequired:        true,
				IsActive:          true,
				Config:            `{"section": "personal"}`,
			},
			{
				LoanTypeVersionID: loanTypeVersion.ID,
				Label:             "Información Financiera",
				Code:              "financial_info",
				Description:       "Datos financieros del solicitante",
				Order:             2,
				IsRequired:        true,
				IsActive:          true,
				Config:            `{"section": "financial"}`,
			},
			{
				LoanTypeVersionID: loanTypeVersion.ID,
				Label:             "Detalles del Préstamo",
				Code:              "loan_details",
				Description:       "Información específica del préstamo solicitado",
				Order:             3,
				IsRequired:        true,
				IsActive:          true,
				Config:            `{"section": "loan"}`,
			},
		}

		// Insertar todos los formularios de una vez
		if err := DB.Create(&forms).Error; err != nil {
			return err
		}
	}

	// Obtener los formularios creados
	var createdForms []models.LoanTypeForm
	DB.Where("loan_type_version_id = ?", loanTypeVersion.ID).Find(&createdForms)

	// Verificar si ya existen inputs
	var existingInputsCount int64
	formIDs := make([]uint, len(createdForms))
	for i, form := range createdForms {
		formIDs[i] = form.ID
	}
	DB.Model(&models.LoanTypeVersionFormInput{}).Where("loan_type_form_id IN ?", formIDs).Count(&existingInputsCount)

	if existingInputsCount == 0 {
		// Crear mapa de formularios por código para fácil acceso
		formMap := make(map[string]uint)
		for _, form := range createdForms {
			formMap[form.Code] = form.ID
		}

		// Preparar todos los inputs de una vez
		allInputs := []models.LoanTypeVersionFormInput{
			// Inputs para Información Personal
			{
				LoanTypeFormID:  formMap["personal_info"],
				Label:           "Nombre Completo",
				Code:            "full_name",
				InputType:       "text",
				Placeholder:     "Ingrese su nombre completo",
				ValidationRules: `{"required": true, "minLength": 2}`,
				Options:         `{}`,
				Order:           1,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
			{
				LoanTypeFormID:  formMap["personal_info"],
				Label:           "Tipo de Documento",
				Code:            "document_type",
				InputType:       "select",
				Placeholder:     "Seleccione el tipo de documento",
				ValidationRules: `{"required": true}`,
				Options:         `["cedula", "pasaporte", "tarjeta_identidad"]`,
				Order:           2,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
			{
				LoanTypeFormID:  formMap["personal_info"],
				Label:           "Número de Documento",
				Code:            "document_number",
				InputType:       "text",
				Placeholder:     "Ingrese su número de documento",
				ValidationRules: `{"required": true, "minLength": 5, "maxLength": 20}`,
				Options:         `{}`,
				Order:           3,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
			{
				LoanTypeFormID:  formMap["personal_info"],
				Label:           "Edad",
				Code:            "age",
				InputType:       "number",
				Placeholder:     "Ingrese su edad",
				ValidationRules: `{"required": true, "min": 18, "max": 75}`,
				Options:         `{}`,
				Order:           4,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
			// Inputs para Información Financiera
			{
				LoanTypeFormID:  formMap["financial_info"],
				Label:           "Ingresos Mensuales",
				Code:            "monthly_income",
				InputType:       "number",
				Placeholder:     "Ingrese sus ingresos mensuales",
				ValidationRules: `{"required": true, "min": 1000000}`,
				Options:         `{}`,
				Order:           1,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
			{
				LoanTypeFormID:  formMap["financial_info"],
				Label:           "Gastos Mensuales",
				Code:            "monthly_expenses",
				InputType:       "number",
				Placeholder:     "Ingrese sus gastos mensuales",
				ValidationRules: `{"required": true, "min": 0}`,
				Options:         `{}`,
				Order:           2,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
			// Inputs para Detalles del Préstamo
			{
				LoanTypeFormID:  formMap["loan_details"],
				Label:           "Monto Solicitado",
				Code:            "requested_amount",
				InputType:       "number",
				Placeholder:     "Ingrese el monto que desea solicitar",
				ValidationRules: `{"required": true, "min": 100000, "max": 10000000}`,
				Options:         `{}`,
				Order:           1,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
			{
				LoanTypeFormID:  formMap["loan_details"],
				Label:           "Propósito del Préstamo",
				Code:            "purpose",
				InputType:       "select",
				Placeholder:     "Seleccione el propósito",
				ValidationRules: `{"required": true}`,
				Options:         `["Educación", "Vivienda", "Vehículo", "Consolidación de deudas", "Negocio", "Gastos médicos", "Otros"]`,
				Order:           2,
				IsRequired:      true,
				IsActive:        true,
				Config:          `{}`,
			},
		}

		// Insertar todos los inputs de una vez
		if err := DB.Create(&allInputs).Error; err != nil {
			return err
		}
	}

	log.Println("Datos iniciales insertados correctamente")
	return nil
}
