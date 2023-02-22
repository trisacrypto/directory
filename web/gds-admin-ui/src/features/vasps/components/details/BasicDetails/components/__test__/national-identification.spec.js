import { faker } from '@faker-js/faker';

import countryCodeEmoji, { getCountryName } from '@/utils/country';
import { render, screen } from '@/utils/test-utils';
import NationalIdentification from '../NationalIdentification';

describe('<NationalIdentification/>', () => {
    it('should render data', () => {
        const data = {
            customer_number: faker.phone.number(),
            country_of_registration: 'US',
            national_identification: {
                national_identifier: '36-40XXXX',
                national_identifier_type: 'NATIONAL_IDENTIFIER_TYPE_CODE_TXID',
                registration_authority: 'IRS',
            },
        };

        render(<NationalIdentification data={data} />);

        expect(screen.getByText(/issued by:/i).firstElementChild.textContent).toBe(
            `${countryCodeEmoji(data.country_of_registration)} (${data.country_of_registration}) by authority ${
                data.national_identification.registration_authority
            }`
        );
        expect(screen.getByText(/national identification type:/i).firstElementChild.textContent).toBe('Tax ID');
        expect(screen.getByText(/country of registration:/i).firstElementChild.textContent).toBe(
            getCountryName(data.country_of_registration)
        );
        expect(screen.getByText(/customer number:/i).firstElementChild.textContent).toBe(data.customer_number);
    });

    it('should render N/A for empty properties value', () => {
        const data = {
            customer_number: '',
            national_identification: {
                country_of_issue: '',
                national_identifier: '',
                national_identifier_type: '',
                registration_authority: '',
            },
        };

        render(<NationalIdentification data={data} />);

        expect(screen.getByText(/issued by:/i).firstElementChild.textContent).toBe('N/A (N/A) by authority N/A');
        expect(screen.getByText(/national identification type:/i).firstElementChild.textContent).toBe('N/A');
        expect(screen.getByText(/country of registration:/i).firstElementChild.textContent).toBe('N/A');
        expect(screen.getByText(/customer number:/i).firstElementChild.textContent).toBe('N/A');
    });

    it('should render N/A if national_identification property is null', () => {
        const data = {
            customer_number: '',
            national_identification: null,
        };

        render(<NationalIdentification data={data} />);

        expect(screen.getByText(/issued by:/i).firstElementChild.textContent).toBe('N/A (N/A) by authority N/A');
        expect(screen.getByText(/national identification type:/i).firstElementChild.textContent).toBe('N/A');
        expect(screen.getByText(/country of registration:/i).firstElementChild.textContent).toBe('N/A');
        expect(screen.getByText(/customer number:/i).firstElementChild.textContent).toBe('N/A'); // this is the case when the data is not available')
    });
});
