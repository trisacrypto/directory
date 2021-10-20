import { render, screen, waitFor } from "@testing-library/react"
import userEvent from '@testing-library/user-event'
import BasicDetails from "../pages/app/details/BasicDetails"
import faker from "faker"
import { Status } from "../constants"

describe("BasicDetailsDropDown", () => {

    it("should render", () => {
        render(<BasicDetails />)
    })

    it("Should not be disable when status is different to pending", () => {

        const mockVaspData = {
            "name": "Opalcliff, Inc.",
            "vasp": {
                "business_category": "BUSINESS_ENTITY",
                "common_name": "trisa.opalcliff.com",
                "contacts": {
                    "administrative": null,
                    "billing": null,
                    "legal": null,
                    "technical": {
                        "email": "trisa@opalcliff.com",
                        "extra": null,
                        "name": "jj",
                        "person": null,
                        "phone": "5838805329"
                    }
                },
                "entity": {
                    "country_of_registration": "CA",
                    "customer_number": "",
                    "geographic_addresses": [
                        {
                            "address_line": [
                                "140 Victory Lane",
                                "Los Gatos",
                                ""
                            ],
                            "address_type": "ADDRESS_TYPE_CODE_BIZZ",
                            "building_name": "",
                            "building_number": "",
                            "country": "CA",
                            "country_sub_division": "",
                            "department": "",
                            "district_name": "",
                            "floor": "",
                            "post_box": "",
                            "post_code": "",
                            "room": "",
                            "street_name": "",
                            "sub_department": "",
                            "town_location_name": "",
                            "town_name": ""
                        }
                    ],
                    "name": {
                        "local_name_identifiers": [
                            {
                                "legal_person_name": "Opalcliff",
                                "legal_person_name_identifier_type": "LEGAL_PERSON_NAME_TYPE_CODE_SHRT"
                            }
                        ],
                        "name_identifiers": [
                            {
                                "legal_person_name": "Opalcliff, Inc.",
                                "legal_person_name_identifier_type": "LEGAL_PERSON_NAME_TYPE_CODE_LEGL"
                            }
                        ],
                        "phonetic_name_identifiers": []
                    },
                    "national_identification": {
                        "country_of_issue": "US",
                        "national_identifier": "5",
                        "national_identifier_type": "NATIONAL_IDENTIFIER_TYPE_CODE_LEIX",
                        "registration_authority": ""
                    }
                },
                "established_on": "275760-08-08",
                "extra": null,
                "first_listed": "2021-07-28T15:44:23Z",
                "id": "30f624c5-2007-4c2a-8770-31e5087d042d",
                "identity_certificate": null,
                "last_updated": "2021-07-28T15:50:06Z",
                "registered_directory": "trisatest.net",
                "service_status": "UNKNOWN",
                "signature": "",
                "signing_certificates": [],
                "trisa_endpoint": "opalcliff.com:443",
                "trixo": {
                    "applicable_regulations": [
                        "FATF Recommendation 16"
                    ],
                    "compliance_threshold": 0,
                    "compliance_threshold_currency": "USD",
                    "conducts_customer_kyc": false,
                    "financial_transfers_permitted": "",
                    "has_required_regulatory_program": "",
                    "kyc_threshold": 0,
                    "kyc_threshold_currency": "USD",
                    "must_comply_travel_rule": false,
                    "must_safeguard_pii": false,
                    "other_jurisdictions": [],
                    "primary_national_jurisdiction": "US",
                    "primary_regulator": "Wife",
                    "safeguards_pii": false
                },
                "vasp_categories": [
                    "P2P"
                ],
                "verification_status": "PENDING_REVIEW",
                "verified_on": "",
                "version": {
                    "pid": "0",
                    "version": "4"
                },
                "website": "http://opalcliff.com"
            },
            "verified_contacts": {
                "technical": "trisa@opalcliff.com"
            },
            "traveler": false
        }
        render(<BasicDetails data={mockVaspData} />)
        const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

        waitFor(() => userEvent.click(dropdownDots))
        const dropdown = screen.getByTestId(/reviewItem/i)

        expect(dropdown).not.toHaveClass('disabled')
    })

    it("Should be disabled when status is pending", () => {
        const mockVaspData = {
            "name": "Opalcliff, Inc.",
            "vasp": {
                "business_category": "BUSINESS_ENTITY",
                "common_name": "trisa.opalcliff.com",
                "contacts": {
                    "administrative": null,
                    "billing": null,
                    "legal": null,
                    "technical": {
                        "email": "trisa@opalcliff.com",
                        "extra": null,
                        "name": "jj",
                        "person": null,
                        "phone": "5838805329"
                    }
                },
                "entity": {
                    "country_of_registration": "CA",
                    "customer_number": "",
                    "geographic_addresses": [
                        {
                            "address_line": [
                                "140 Victory Lane",
                                "Los Gatos",
                                ""
                            ],
                            "address_type": "ADDRESS_TYPE_CODE_BIZZ",
                            "building_name": "",
                            "building_number": "",
                            "country": "CA",
                            "country_sub_division": "",
                            "department": "",
                            "district_name": "",
                            "floor": "",
                            "post_box": "",
                            "post_code": "",
                            "room": "",
                            "street_name": "",
                            "sub_department": "",
                            "town_location_name": "",
                            "town_name": ""
                        }
                    ],
                    "name": {
                        "local_name_identifiers": [
                            {
                                "legal_person_name": "Opalcliff",
                                "legal_person_name_identifier_type": "LEGAL_PERSON_NAME_TYPE_CODE_SHRT"
                            }
                        ],
                        "name_identifiers": [
                            {
                                "legal_person_name": "Opalcliff, Inc.",
                                "legal_person_name_identifier_type": "LEGAL_PERSON_NAME_TYPE_CODE_LEGL"
                            }
                        ],
                        "phonetic_name_identifiers": []
                    },
                    "national_identification": {
                        "country_of_issue": "US",
                        "national_identifier": "5",
                        "national_identifier_type": "NATIONAL_IDENTIFIER_TYPE_CODE_LEIX",
                        "registration_authority": ""
                    }
                },
                "established_on": "275760-08-08",
                "extra": null,
                "first_listed": "2021-07-28T15:44:23Z",
                "id": "30f624c5-2007-4c2a-8770-31e5087d042d",
                "identity_certificate": null,
                "last_updated": "2021-07-28T15:50:06Z",
                "registered_directory": "trisatest.net",
                "service_status": "UNKNOWN",
                "signature": "",
                "signing_certificates": [],
                "trisa_endpoint": "opalcliff.com:443",
                "trixo": {
                    "applicable_regulations": [
                        "FATF Recommendation 16"
                    ],
                    "compliance_threshold": 0,
                    "compliance_threshold_currency": "USD",
                    "conducts_customer_kyc": false,
                    "financial_transfers_permitted": "",
                    "has_required_regulatory_program": "",
                    "kyc_threshold": 0,
                    "kyc_threshold_currency": "USD",
                    "must_comply_travel_rule": false,
                    "must_safeguard_pii": false,
                    "other_jurisdictions": [],
                    "primary_national_jurisdiction": "US",
                    "primary_regulator": "Wife",
                    "safeguards_pii": false
                },
                "vasp_categories": [
                    "P2P"
                ],
                "verification_status": faker.random.arrayElement([Status.APPEALED, Status.EMAIL_VERIFIED, Status.REJECTED, Status.SUBMITTED, Status.ERRORED, Status.ISSUING_CERTIFICATE]),
                "verified_on": "",
                "version": {
                    "pid": "0",
                    "version": "4"
                },
                "website": "http://opalcliff.com"
            },
            "verified_contacts": {
                "technical": "trisa@opalcliff.com"
            },
            "traveler": false
        }

        render(<BasicDetails data={mockVaspData} />)
        const dropdownDots = screen.getByTestId(/dripicons-dots-3/i)

        waitFor(() => userEvent.click(dropdownDots))
        const dropdown = screen.getByTestId(/reviewItem/i)

        expect(dropdown).toHaveClass('disabled')
    })
})