const normalizeFlatArrays = (data) => data && data.map(d => ({ name: d }))

export const getTrixoFormInitialValues = (data = []) => {
    const defaultValues = {
        applicable_regulations: [{
            name: "FATF Recommendation 16"
        }],
        compliance_threshold: 0,
        compliance_threshold_currency: "",
        conducts_customer_kyc: false,
        financial_transfers_permitted: "",
        has_required_regulatory_program: "",
        kyc_threshold: 0,
        kyc_threshold_currency: "",
        must_comply_travel_rule: false,
        must_safeguard_pii: false,
        other_jurisdictions: [{
            country: "",
            regulator_name: ""
        }],
        primary_national_jurisdiction: "",
        primary_regulator: "",
        safeguards_pii: false
    }

    const applicableRegulations = normalizeFlatArrays(data.applicable_regulations)

    const result = {
        ...Object.assign(defaultValues, data),
        applicable_regulations: applicableRegulations,
    }

    return result
}

export const getBusinessInfosFormInitialValues = (data) => ({
    website: data.vasp.website || "",
    established_on: data.vasp.established_on || "",
    vasp_categories: data.vasp.vasp_categories || [],
    business_category: data.vasp.business_category || ""
})

export const getTrisaImplementationDetailsInitialValue = (data) => ({
    common_name: data.vasp.common_name,
    trisa_endpoint: data.vasp.trisa_endpoint
})

export const getIvms101RecordInitialValues = (data) => {
    const defaultValues = {
        name: {
            name_identifiers: [
                {
                    legal_person_name: "",
                    legal_person_name_identifier_type: "0"
                }
            ],
            local_name_identifiers: [],
            phonetic_name_identifiers: []
        },
        geographic_addresses: [
            {
                address_type: 2,
                address_line: [
                    "",
                    "",
                    ""
                ],
                country: ""
            }
        ],
        customer_number: "",
        national_identification: {
            national_identifier: "",
            national_identifier_type: 0,
            country_of_issue: "",
            registration_authority: ""
        },
        country_of_registration: ""
    }

    return {
        ...Object.assign(defaultValues, data),
    }
}

export const getContactInitialValues = (data) => {
    const initialValue = {
        email: "",
        extra: null,
        name: "",
        person: null,
        phone: ""
    }

    return Object.assign(initialValue, data)
}