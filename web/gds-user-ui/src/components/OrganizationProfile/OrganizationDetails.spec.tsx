import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import OrganizationalDetail from './OrganizationalDetail';

describe('<OrganizationDetails />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render properties correctly', () => {
    const data = {
      entity: {
        country_of_registration: 'AS',
        name: {
          name_identifiers: [
            {
              legal_person_name: 'VASP Holding',
              legal_person_name_identifier_type: 'LEGAL_PERSON_NAME_TYPE_CODE_LEGL'
            },
            {
              legal_person_name: '',
              legal_person_name_identifier_type: ''
            }
          ],
          local_name_identifiers: [],
          phonetic_name_identifiers: []
        },
        geographic_addresses: [
          {
            address_type: 'ADDRESS_TYPE_CODE_BIZZ',
            address_line: ['50 Oakland Ave', '602'],
            country: 'AO',
            city: 'Mansfield',
            state: 'Ohio',
            postal_code: '06268'
          }
        ],
        national_identification: {
          national_identifier_type: 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX',
          country_of_issue: '',
          registration_authority: 'RA777777',
          national_identifier: '3'
        }
      },
      contacts: {
        administrative: {
          name: 'mulary',
          email: 'mulary@mailinator.com'
        },
        technical: {
          name: 'refum',
          email: 'refum@mailinator.com'
        },
        billing: {
          name: 'zuca',
          email: 'zuca@mailinator.com'
        },
        legal: {
          name: 'refus',
          email: 'refus@mailinator.com',
          phone: '+264940993994'
        }
      },
      trisa_endpoint_testnet: {
        trisa_endpoint: '',
        common_name: 'testnet.sixelen.me',
        endpoint: 'testnet.sixelen.me:443'
      },
      trisa_endpoint_mainnet: {
        trisa_endpoint: '',
        common_name: 'trisa.sixelen.me',
        endpoint: 'trisa.sixelen.me:443'
      },
      website: 'https://www.sixelen.me',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['P2P'],
      established_on: '2022-07-07',
      organization_name: 'VASP Holding',
      trixo: {
        primary_national_jurisdiction: 'AS',
        primary_regulator: '',
        other_jurisdictions: [],
        financial_transfers_permitted: 'no',
        has_required_regulatory_program: 'no',
        conducts_customer_kyc: false,
        kyc_threshold: 0,
        kyc_threshold_currency: 'USD',
        must_comply_travel_rule: false,
        applicable_regulations: [
          {
            name: 'FATF Recommendation 16'
          }
        ],
        compliance_threshold: 3000,
        compliance_threshold_currency: 'USD',
        must_safeguard_pii: false,
        safeguards_pii: false
      }
    };
    render(<OrganizationalDetail data={data} />);

    expect(screen.getByTestId('legal_person_name').textContent).toBe('VASP Holding');
    expect(screen.getByTestId('business_category').textContent).toBe('Business Entity');

    expect(screen.getByTestId('addressLine').textContent).toBe('50 Oakland Ave 602 AO');
    expect(screen.getByTestId('vasp_categories').textContent).toBe('Person-to-Person Exchange');
    expect(screen.getByTestId('country_of_registration').textContent).toBe('American Samoa');
    expect(screen.getByTestId('national_identifier_type').textContent).toBe(
      'Legal Entity Identifier (LEI)'
    );

    const contactNames = screen
      .getAllByTestId('contact__name')
      .map((contact) => contact.textContent);
    expect(contactNames).toEqual(['refus', 'refum', 'mulary', 'zuca']);

    const contactEmails = screen.getAllByTestId('contact__email').map((email) => email.textContent);
    expect(contactEmails).toEqual([
      'refus@mailinator.com',
      'refum@mailinator.com',
      'mulary@mailinator.com',
      'zuca@mailinator.com'
    ]);

    const contactPhones = screen.getAllByTestId('contact__phone').map((phone) => phone.textContent);
    expect(contactPhones).toEqual(['+264940993994']);
  });

  it('should display N/A when there no data', () => {
    const data = {};
    render(<OrganizationalDetail data={data} />);

    expect(screen.getByTestId('legal_person_name').textContent).toBe('N/A');
    expect(screen.getByTestId('business_category').textContent).toBe('N/A');

    expect(screen.getByTestId('vasp_categories').textContent).toBe('N/A');
    expect(screen.getByTestId('country_of_registration').textContent).toBe('N/A');
    expect(screen.getByTestId('national_identifier_type').textContent).toBe('N/A');

    const contactNames = screen.getAllByTestId('contacts').map((contact) => contact.textContent);
    expect(contactNames).toEqual(['N/A', 'N/A', 'N/A', 'N/A']);
  });
});
