import { t } from '@lingui/macro';
import * as yup from 'yup';
import { setupI18n } from '@lingui/core';

const _i18n = setupI18n();

const trisaEndpointPattern = /^([a-zA-Z0-9.-]+):((?!(0))[0-9]+)$/;
const commonNameRegex =
    /^([a-z0-9]+([-a-z0-9]*[a-z0-9]+)?\.){0,}([a-z0-9]+([-a-z0-9]*[a-z0-9]+)?){1,63}(\.[a-z0-9]{2,7})+$/;

export const validationSchema = [
    yup.object().shape({
        website: yup
            .string()
            .url()
            .trim()
            .required(_i18n._(t`Website is a required field`)),
        established_on: yup
            .date()
            .nullable()
            .transform((curr, orig) => (orig === '' ? null : curr))
            .required(_i18n._(t`Invalid date`))
            .test('is-invalidate-date', _i18n._(t`Invalid date / year must be 4 digit`), (value) => {
                if (value) {
                    const getYear = value.getFullYear();
                    if (getYear.toString().length !== 4) {
                        return false;
                    } else {
                        return true;
                    }
                }
                return false;
            })
            .required(),
        organization_name: yup
            .string()
            .trim()
            .required(_i18n._(t`Organization name is required`)),
        business_category: yup.string().nullable(true),
        vasp_categories: yup.array().of(yup.string()).nullable(true)
    }),
    yup.object().shape({
        entity: yup.object().shape({
            country_of_registration: yup
                .string()
                .required(_i18n._(t`Country of registration is required`)),
            name: yup.object().shape({
                name_identifiers: yup.array(
                    yup.object().shape({
                        legal_person_name: yup
                            .string()
                            .test(
                                'notEmptyIfIdentifierTypeExist',
                                _i18n._(t`Legal name is required`),
                                (value, ctx): any => {
                                    return !(ctx.parent.legal_person_name_identifier_type && !value);
                                }
                            ),
                        legal_person_name_identifier_type: yup.string().when('legal_person_name', {
                            is: (value: string) => !!value,
                            then: yup.string().required(_i18n._(t`Name Identifier Type is required`))
                        })
                    })
                ),
                local_name_identifiers: yup.array(
                    yup.object().shape({
                        legal_person_name: yup
                            .string()
                            .test(
                                'notEmptyIfIdentifierTypeExist',
                                _i18n._(t`Legal name is required`),
                                (value, ctx): any => {
                                    return !(ctx.parent.legal_person_name_identifier_type && !value);
                                }
                            ),
                        legal_person_name_identifier_type: yup.string().when('legal_person_name', {
                            is: (value: string) => !!value,
                            then: yup.string().required(_i18n._(t`Name Identifier Type is required`))
                        })
                    })
                ),
                phonetic_name_identifiers: yup.array(
                    yup.object().shape({
                        legal_person_name: yup
                            .string()
                            .test(
                                'notEmptyIfIdentifierTypeExist',
                                _i18n._(t`Legal name is required`),
                                (value, ctx): any => {
                                    return !(ctx.parent.legal_person_name_identifier_type && !value);
                                }
                            ),
                        legal_person_name_identifier_type: yup.string().when('legal_person_name', {
                            is: (value: string) => !!value,
                            then: yup.string().required(_i18n._(t`Name Identifier Type is required`))
                        })
                    })
                )
            }),
            geographic_addresses: yup.array().of(
                yup.object().shape({
                    address_line: yup.array(),
                    'address_line[0]': yup
                        .string()
                        .test('test-0', 'addresse line 0', (value: any, ctx: any): any => {
                            return ctx && ctx.parent && ctx.parent.address_line[0];
                        }),
                    'address_line[2]': yup
                        .string()
                        .test('test-0', 'addresse line 0', (value: any, ctx: any): any => {
                            return ctx && ctx.parent && ctx.parent.address_line[2];
                        }),
                    country: yup.string().required(),
                    postal_code: yup.string().required(),
                    state: yup.string().required(),
                    address_type: yup.string().required(),
                })
            ),
            national_identification: yup.object().shape({
                national_identifier: yup.string().required(),
                national_identifier_type: yup.string(),
                country_of_issue: yup.string(),
                registration_authority: yup
                    .string()
                    .test(
                        'registrationAuthority',
                        _i18n._(t`Registration Authority cannot be left empty`),
                        (value, ctx) => {
                            if (
                                ctx.parent.national_identifier_type !== 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX' &&
                                !value
                            ) {
                                return false;
                            }

                            return true;
                        }
                    )
            })
        })
    }),
    yup.object().shape({
        contacts: yup.object().shape({
            administrative: yup.object().shape({
                name: yup.string(),
                email: yup.string().email(_i18n._(t`Email is not valid`)),
                phone: yup.string()
            }),
            technical: yup
                .object()
                .shape({
                    name: yup.string().required(),
                    email: yup
                        .string()
                        .email(_i18n._(t`Email is not valid`))
                        .required(_i18n._(t`Email is required`)),
                    phone: yup.string()
                })
                .required(),
            billing: yup.object().shape({
                name: yup.string(),
                email: yup.string().email(_i18n._(t`Email is not valid`)),
                phone: yup.string()
            }),
            legal: yup
                .object()
                .shape({
                    name: yup.string().required(),
                    email: yup.string().email('Email is not valid').required('Email is required'),
                    phone: yup
                        .string()
                        .required(
                            'A business phone number is required to complete physical verification for MainNet registration. Please provide a phone number where the Legal/ Compliance contact can be contacted.'
                        )
                })
                .required()
        })
    }),
    yup.object().shape({
        trisa_endpoint: yup.string().trim(),
        trisa_endpoint_testnet: yup.object().shape({
            endpoint: yup.string().matches(trisaEndpointPattern, _i18n._(t`TRISA endpoint is not valid`)),
            common_name: yup
                .string()
                .matches(
                    commonNameRegex,
                    _i18n._(
                        t`Common name should not contain special characters, no spaces and must have a dot(.) in it`
                    )
                )
        }),
        trisa_endpoint_mainnet: yup.object().shape({
            endpoint: yup
                .string()
                .test(
                    'uniqueMainetEndpoint',
                    _i18n._(t`TestNet and MainNet endpoints should not be the same`),
                    (value, ctx: any): any => {
                        return ctx.from[1].value.trisa_endpoint_testnet.endpoint !== value;
                    }
                )
                .matches(trisaEndpointPattern, _i18n._(t`TRISA endpoint is not valid`)),
            common_name: yup
                .string()
                .matches(
                    commonNameRegex,
                    _i18n._(
                        t`Common name should not contain special characters, no spaces and must have a dot(.) in it`
                    )
                )
        })
    }),
    yup.object().shape({
        trixo: yup.object().shape({
            primary_national_jurisdiction: yup.string(),
            primary_regulator: yup.string(),
            other_jurisdictions: yup.array().of(
                yup.object().shape({
                    country: yup.string(),
                    regulator_name: yup.string()
                })
            ),
            financial_transfers_permitted: yup.string().oneOf(['no', 'yes', 'partial']).default('no'),
            has_required_regulatory_program: yup.string().oneOf(['no', 'yes', 'partial']).default('no'),
            conducts_customer_kyc: yup.boolean().default(false),
            kyc_threshold: yup.number(),
            kyc_threshold_currency: yup.string(),
            must_comply_travel_rule: yup.boolean(),
            applicable_regulations: yup
                .array()
                .of(
                    yup.object().shape({
                        name: yup.string()
                    })
                )
                .transform((value, originalValue) => {
                    if (originalValue) {
                        return originalValue.filter((item: any) => item.name.length > 0);
                    }
                    return value;

                    // remove empty items
                }),
            compliance_threshold: yup.number(),
            compliance_threshold_currency: yup.string(),
            must_safeguard_pii: yup.boolean().default(false),
            safeguards_pii: yup.boolean().default(false)
        })
    })
];
