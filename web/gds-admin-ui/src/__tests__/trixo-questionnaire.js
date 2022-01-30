import TrixoQuestionnaire from "pages/app/details/TrixoQuestionnaire"
import countryCodeEmoji, { isoCountries } from "utils/country";
import { render, screen } from "utils/test-utils"
// eslint-disable-next-line jest/no-mocks-import
import generateOtherJuridictions from "../__mocks__/other-juridictions";


describe("<TrixoQuestionnaire />", () => {

    describe('financial_transfers_permitted', () => {

        it("should be permitted", () => {
            const data = { financial_transfers_permitted: "yes" }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("financial-transfers-permitted").textContent).toBe("Organization is permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates")
        })

        it("should not be permitted", () => {
            const data = { financial_transfers_permitted: "no" }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("financial-transfers-permitted").textContent).toBe("Organization is not permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates")
        })

    });

    it("primary juridiction should be well formatted", () => {
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
        const primaryJuridictionEl = screen.getByTestId("primary-juridiction")

        expect(primaryJuridictionEl.textContent).toBe(`${countryCodeEmoji(data.primary_national_jurisdiction)} ${isoCountries[data.primary_national_jurisdiction]} regulated by ${data.primary_regulator}`)

    })


    it("other juridictions should be well formatted", () => {
        const data = {
            applicable_regulations: [
                "FATF Recommendation 16"
            ],
            other_jurisdictions: [
                ...generateOtherJuridictions()
            ],
            primary_national_jurisdiction: "US",
        }

        render(<TrixoQuestionnaire data={data} />)
        const otherJuridictionsEl = screen.getAllByTestId("other-juridiction")
        const otherJuridictionContent = otherJuridictionsEl.map(juridiction => {
            return juridiction.textContent
        })

        expect(otherJuridictionContent[0]).toBe(`${countryCodeEmoji(data.other_jurisdictions[0].country)} ${isoCountries[data.other_jurisdictions[0].country]} regulated by ${data.other_jurisdictions[0].regulator_name}`)
        expect(otherJuridictionContent[1]).toBe(`${countryCodeEmoji(data.other_jurisdictions[1].country)} ${isoCountries[data.other_jurisdictions[1].country]} regulated by ${data.other_jurisdictions[1].regulator_name}`)

    })

    describe('must_comply_travel_rule', () => {

        it("organisation must comply", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("must-comply-travel-rule").textContent).toBe("Organization must comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered.")
        })

        it("organisation must not comply", () => {
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

        it("organisation does programme", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true,
                has_required_regulatory_program: "yes"
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("has-required-regulatory-program").textContent).toBe("Organization does programme that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.")
        })

        it("organisation does not programme", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: false,
                must_safeguard_pii: true,
                has_required_regulatory_program: "no"
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("has-required-regulatory-program").textContent).toBe("Organization does not programme that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.")
        })

        it("organisation has a partially implemented programme", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: false,
                must_safeguard_pii: true,
                has_required_regulatory_program: "partial"
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("has-required-regulatory-program").textContent).toBe("Organization has a partially implemented programme that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.")
        })

    });

    describe('conducts_customer_kyc', () => {

        it("organisation does conduct", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true,
                conducts_customer_kyc: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("conducts-customer-kyc").textContent).toBe("Organization does conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers.")
        })

        it("organisation does not conduct", () => {
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

        it("organisation must safe guard", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                must_safeguard_pii: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("must-safeguard-pii").textContent).toBe("Organization must safeguard PII by law.")
        })

        it("organisation is not required to safe guard", () => {
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

        it("organisation does secure", () => {
            const data = {
                financial_transfers_permitted: "yes"
                , must_comply_travel_rule: true,
                safeguards_pii: true
            }

            render(<TrixoQuestionnaire data={data} />)
            expect(screen.getByTestId("safeguards-pii").textContent).toBe("Organization does secure and protect PII, including PII received from other VASPs under the Travel Rule.")
        })

        it("organisation is not required to safe guard", () => {
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