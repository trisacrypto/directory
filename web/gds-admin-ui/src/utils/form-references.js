import { NAME_IDENTIFIER_TYPE_CODE } from "constants/basic-details"
import { NATIONAL_IDENTIFIER_TYPE_CODE } from "constants/national-identification"

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

function getLegalPersonNameIdentifierTypeCode(nameIdentifiers = []) {
    return nameIdentifiers.map(name => {
        const code = NAME_IDENTIFIER_TYPE_CODE[name.legal_person_name_identifier_type] || "0"
        return {
            ...name,
            legal_person_name_identifier_type: code
        }
    })
}

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
    const initialValues = Object.assign(defaultValues, data)

    return {
        ...Object.assign(defaultValues, data),
        name: {
            name_identifiers: getLegalPersonNameIdentifierTypeCode(initialValues.name.name_identifiers),
            local_name_identifiers: getLegalPersonNameIdentifierTypeCode(initialValues.name.local_name_identifiers),
            phonetic_name_identifiers: getLegalPersonNameIdentifierTypeCode(initialValues.name.phonetic_name_identifiers),
        },
        national_identification: {
            ...initialValues.national_identification,
            national_identifier_type: NATIONAL_IDENTIFIER_TYPE_CODE[initialValues.national_identification.national_identifier_type]
        }
    }
}