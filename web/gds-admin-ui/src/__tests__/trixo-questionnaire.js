import TrixoQuestionnaire from "pages/app/details/TrixoQuestionnaire"
import countryCodeEmoji, { isoCountries } from "utils/country";
import { render, screen } from "utils/test-utils"
// eslint-disable-next-line jest/no-mocks-import
import generateOtherJurisdictions from "../__mocks__/other-jurisdictions";


describe("<TrixoQuestionnaire />", () => {

    describe('financial_transfers_permitted', () => {

        it("should be permitted", () => {
            const data = { financial_transfers_permitted: "yes" }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("financial-transfers-permitted").textContent).toBe("Organization is permitted to send and/or receive transfers of virtual assets in the jurisdiction(s) in which it operates")
        })

        it("should not be permitted", () => {
            const data = { financial_transfers_permitted: "no" }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("financial-transfers-permitted").textContent).toBe("Organization is not permitted to send and/or receive transfers of virtual assets in the jurisdiction(s) in which it operates")
        })

    });

    it("primary jurisdiction should be well formatted", () => {
        const data = {
            applicable_regulations: [
                "FATF Recommendation 16"
            ],
            compliance_threshold: 3000,
            compliance_threshold_currency: "USD",
            conducts_customer_kyc: true,
            financial_transfers_permitted: "no",
            has_required_regulatory_program: "yes",
            kyc_threshold: 0,
            kyc_threshold_currency: "USD",
            must_comply_travel_rule: true,
            must_safeguard_pii: true,
            other_jurisdictions: [],
            primary_regulator: "FinCEN",
            safeguards_pii: true,
            primary_national_jurisdiction: "US",
        }


        render(<TrixoQuestionnaire data={data} />)
        const primaryJurisdictionEl = screen.getByTestId("primary-jurisdiction")

        expect(primaryJurisdictionEl.textContent).toBe(`${countryCodeEmoji(data.primary_national_jurisdiction)} ${isoCountries[data.primary_national_jurisdiction]} regulated by ${data.primary_regulator}`)

    })


    it("other juridisctions should be well formatted", () => {
        const data = {
            applicable_regulations: [
                "FATF Recommendation 16"
            ],
            other_jurisdictions: [
                ...generateOtherJurisdictions()
            ],
            primary_national_jurisdiction: "US",
        }

        render(<TrixoQuestionnaire data={data} />)
        const otherJurisdictionsEl = screen.getAllByTestId("other-jurisdiction")
        const otherJurisdictionContent = otherJurisdictionsEl.map(jurisdiction => {
            return jurisdiction.textContent
        })

        expect(otherJurisdictionContent[0]).toBe(`${countryCodeEmoji(data.other_jurisdictions[0].country)} ${isoCountries[data.other_jurisdictions[0].country]} regulated by ${data.other_jurisdictions[0].regulator_name}`)
        expect(otherJurisdictionContent[1]).toBe(`${countryCodeEmoji(data.other_jurisdictions[1].country)} ${isoCountries[data.other_jurisdictions[1].country]} regulated by ${data.other_jurisdictions[1].regulator_name}`)

    })

    describe('must_comply_travel_rule', () => {

        it("organization must comply", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("must-comply-travel-rule").textContent).toBe("Organization must comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered.")
        })

        it("organization must not comply", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: false,
                must_safeguard_pii: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("must-comply-travel-rule").textContent).toBe("Organization must not comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered.")
        })

    });

    it('minimum compliance threshold should be well formated', () => {
        const data = {
            financial_transfers_permitted: "yes"
            , must_comply_travel_rule: false,
            must_safeguard_pii: true,
            compliance_threshold: 10000
        }

        render(<TrixoQuestionnaire data={data} />)
        expect(screen.getByTestId('compliance_threshold_currency').textContent).toBe(`$10,000.00`)
    })

    describe('has_required_regulatory_program', () => {

        it("organization does program", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true,
                has_required_regulatory_program: "yes"
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("has-required-regulatory-program").textContent).toBe("Organization does program that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.")
        })

        it("organization does not program", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: false,
                must_safeguard_pii: true,
                has_required_regulatory_program: "no"
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("has-required-regulatory-program").textContent).toBe("Organization does not program that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.")
        })

        it("organization has a partially implemented program", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: false,
                must_safeguard_pii: true,
                has_required_regulatory_program: "partial"
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("has-required-regulatory-program").textContent).toBe("Organization has a partially implemented program that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.")
        })

    });

    describe('conducts_customer_kyc', () => {

        it("organization does conduct", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true,
                conducts_customer_kyc: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("conducts-customer-kyc").textContent).toBe("Organization does conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers.")
        })

        it("organization does not conduct", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true,
                conducts_customer_kyc: false
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("conducts-customer-kyc").textContent).toBe("Organization does not conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers.")
        })
    });

    describe('must_safeguard_pii', () => {

        it("organization must safe guard", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("must-safeguard-pii").textContent).toBe("Organization must safeguard PII by law.")
        })

        it("organization is not required to safe guard", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: false,
                conducts_customer_kyc: false
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("must-safeguard-pii").textContent).toBe("Organization is not required to safeguard PII by law.")
        })
    });

    describe('safeguards_pii', () => {

        it("organization does secure", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                safeguards_pii: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("safeguards-pii").textContent).toBe("Organization does secure and protect PII, including PII received from other VASPs under the Travel Rule.")
        })

        it("organization is not required to safe guard", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                safeguards_pii: false,
                conducts_customer_kyc: false
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("safeguards-pii").textContent).toBe("Organization does not secure and protect PII, including PII received from other VASPs under the Travel Rule.")
        })
    });

})