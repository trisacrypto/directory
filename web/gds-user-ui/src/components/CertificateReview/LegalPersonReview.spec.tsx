import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render } from 'utils/test-utils';
import LegalPersonReview from './LegalPersonReview';
// import * as reactRedux from 'react-redux';

// const useSelectorMock = jest.spyOn(reactRedux, 'useSelector');
// const defaultValues = {
//   entity: {
//     country_of_registration: 'AF',
//     name: {
//       name_identifiers: [
//         {
//           legal_person_name: 'KYC',
//           legal_person_name_identifier_type: 'LEGAL_PERSON_NAME_TYPE_CODE_LEGL'
//         },
//         { legal_person_name: '', legal_person_name_identifier_type: '' }
//       ],
//       local_name_identifiers: [],
//       phonetic_name_identifiers: []
//     },
//     geographic_addresses: [
//       {
//         address_type: 'ADDRESS_TYPE_CODE_BIZZ',
//         address_line: ['a', 'Address 2', 'Address 3'],
//         country: 'AO'
//       },
//       {
//         address_type: 'ADDRESS_TYPE_CODE_HOME',
//         address_line: ['705 Ryan Street', 'Sylvania', 'OH 43560'],
//         country: 'US'
//       },
//       {
//         address_type: 'ADDRESS_TYPE_CODE_HOME',
//         address_line: ['11 Garfield St.', 'Libertyville', 'IL 60048'],
//         country: 'US'
//       }
//     ],
//     national_identification: {
//       national_identifier_type: 'NATIONAL_IDENTIFIER_TYPE_CODE_TXID',
//       country_of_issue: '',
//       registration_authority: 'RA777777',
//       national_identifier: '2'
//     }
//   },
//   website: 'http://kyc.com',
//   business_category: 'GOVERNMENT_ENTITY',
//   vasp_categories: ['P2P', 'Kiosk'],
//   established_on: '2022-04-22',
//   organization_name: 'KYC'
// };

describe('<LegalPersonReview />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should match snapshot', () => {
    // useSelectorMock.mockReturnValue({
    //   ...defaultValues
    // });
    const { container } = render(<LegalPersonReview />);

    expect(container).toMatchSnapshot();
  });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
