import faker from "faker"
import NationalIdentification from "pages/app/details/BasicDetails/components/NationalIdentification"
import countryCodeEmoji, { getCountryName } from "utils/country"
import { render, screen } from "utils/test-utils"

describe("<NationalIdentification/>", () => {


    it("should render data", () => {
        const data = {
            customer_number: faker.phone.phoneNumberFormat(),
            national_identification: {
                country_of_issue: faker.address.countryCode(),
                national_identifier: "36-40XXXX",
                national_identifier_type: "NATIONAL_IDENTIFIER_TYPE_CODE_TXID",
                registration_authority: "IRS"
            }
        }

        render(<NationalIdentification data={data} />)

        expect(screen.getByText(/issued by:/i).firstElementChild.textContent).toBe(`${countryCodeEmoji(data.national_identification.country_of_issue)} (${data.national_identification.country_of_issue}) by authority ${data.national_identification.registration_authority}`)
        expect(screen.getByText(/national identification type:/i).firstElementChild.textContent).toBe('Tax ID')
        expect(screen.getByText(/country of registration:/i).firstElementChild.textContent).toBe(getCountryName(data.national_identification.country_of_issue))
        expect(screen.getByText(/customer number:/i).firstElementChild.textContent).toBe(data.customer_number)

    })

    it("should render N/A for empty properties value", () => {
        const data = {
            customer_number: "",
            national_identification: {
                country_of_issue: "",
                national_identifier: "",
                national_identifier_type: "",
                registration_authority: ""
            }
        }

        render(<NationalIdentification data={data} />)

        expect(screen.getByText(/issued by:/i).firstElementChild.textContent).toBe("N/A (N/A) by authority N/A")
        expect(screen.getByText(/national identification type:/i).firstElementChild.textContent).toBe("N/A")
        expect(screen.getByText(/country of registration:/i).firstElementChild.textContent).toBe("N/A")
        expect(screen.getByText(/customer number:/i).firstElementChild.textContent).toBe("N/A")
    })


    it('should render N/A if national_identification property is null', () => {
        const data = {
            customer_number: "",
            national_identification: null
        }

        render(<NationalIdentification data={data} />)

        expect(screen.getByText(/issued by:/i).firstElementChild.textContent).toBe("N/A (N/A) by authority N/A")
        expect(screen.getByText(/national identification type:/i).firstElementChild.textContent).toBe("N/A")
        expect(screen.getByText(/country of registration:/i).firstElementChild.textContent).toBe("N/A")
        expect(screen.getByText(/customer number:/i).firstElementChild.textContent).toBe("N/A")  // this is the case when the data is not available')
    })
})