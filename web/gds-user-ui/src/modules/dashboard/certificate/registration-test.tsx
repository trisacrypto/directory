// /* eslint-disable require-await */
// import userEvent from '@testing-library/user-event';
// import { dynamicActivate } from 'utils/i18nLoaderHelper';
// import { act, render, screen, waitFor } from 'utils/test-utils';
// import Certificate from './registration';

// const certificateInitialValue = {
//   entity: {
//     country_of_registration: '',
//     name: {
//       name_identifiers: [
//         {
//           legal_person_name: '',
//           legal_person_name_identifier_type: 'LEGAL_PERSON_NAME_TYPE_CODE_LEGL'
//         }
//       ],
//       local_name_identifiers: [],
//       phonetic_name_identifiers: []
//     },
//     geographic_addresses: [{ address_type: '', address_line: ['', '', ''], country: '' }],
//     national_identification: {
//       national_identifier_type: 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX',
//       country_of_issue: '',
//       registration_authority: ''
//     }
//   },
//   contacts: {
//     administrative: { name: '', email: '', phone: '' },
//     technical: { name: '', email: '', phone: '' },
//     billing: {},
//     legal: {}
//   },
//   testnet: { endpoint: '', common_name: '' },
//   mainnet: { endpoint: '', common_name: '' },
//   website: '',
//   business_category: '',
//   vasp_categories: [],
//   established_on: '',
//   organization_name: '',
//   trixo: {
//     primary_national_jurisdiction: '',
//     primary_regulator: '',
//     other_jurisdictions: [],
//     financial_transfers_permitted: 'no',
//     has_required_regulatory_program: 'no',
//     conducts_customer_kyc: false,
//     kyc_threshold: 0,
//     kyc_threshold_currency: 'USD',
//     must_comply_travel_rule: false,
//     applicable_regulations: [{ name: 'FATF Recommendation 16' }],
//     compliance_threshold: 3000,
//     compliance_threshold_currency: 'USD',
//     must_safeguard_pii: false,
//     safeguards_pii: false
//   }
// };

// describe('<Certificate />', () => {
//   beforeAll(() => {
//     act(() => {
//       dynamicActivate('en');
//     });
//   });
//   describe('<BasicDetails />', () => {
//     it('should validate required field', async () => {
//       const initialState = {
//         stepper: {
//           currentStep: 1,
//           steps: [
//             { key: 1, status: 'complete', data: {} },
//             { key: 2, status: 'progress' }
//           ],
//           lastStep: null,
//           hasReachSubmitStep: false
//         }
//       };
//       const basicDetailsValidationMessages = [
//         'Organization name is required',
//         'Website is a required field',
//         'Invalid date / year must be 4 digit'
//       ];

//       await act(async () => {
//         render(<Certificate />, { preloadedState: initialState });
//       });

//       const submitButton = screen.getByRole('button', { name: /save & next/i });
//       userEvent.click(submitButton);

//       await waitFor(() => {
//         const errorMessages = screen
//           .getAllByTestId('error-message')
//           .map((error) => error.textContent);
//         expect(basicDetailsValidationMessages).toEqual(errorMessages);
//       });
//     });

//     describe('website', () => {
//       it('website field should have valid url', async () => {
//         await act(async () => {
//           render(<Certificate />);
//         });

//         const websiteField = screen.getByRole('textbox', { name: /website/i });
//         userEvent.type(websiteField, 'Rotational');

//         const submitButton = screen.getByRole('button', { name: /save & next/i });
//         userEvent.click(submitButton);

//         await waitFor(() => {
//           const errorMessage = screen.getByText(/website must be a valid url/i);
//           expect(errorMessage).toBeInTheDocument();
//           expect(websiteField).toHaveAttribute('aria-invalid', 'true');
//         });
//       });

//       it('website field should have valid url', async () => {
//         await act(async () => {
//           render(<Certificate />);
//         });

//         const websiteField = screen.getByRole('textbox', { name: /website/i });
//         userEvent.type(websiteField, 'http://www.rotational.io');

//         const submitButton = screen.getByRole('button', { name: /save & next/i });
//         userEvent.click(submitButton);

//         expect(websiteField).not.toHaveAttribute('aria-invalid');
//       });
//     });

//     describe('Date of Incorporation / Establishment', () => {
//       beforeAll(() => {
//         act(() => {
//           dynamicActivate('en');
//         });
//       });

//       it('date of corporation field should have valid date', async () => {
//         await act(async () => {
//           render(<Certificate />);
//         });

//         const establishedOnField = screen.getByLabelText(/Date of Incorporation \/ Establishment/i);
//         userEvent.type(establishedOnField, '12/12/22222');

//         const submitButton = screen.getByRole('button', { name: /save & next/i });
//         userEvent.click(submitButton);

//         await waitFor(() => {
//           const errorMessages = screen
//             .getAllByTestId('error-message')
//             .map((err) => err.textContent);

//           expect(errorMessages).toContain('Invalid date / year must be 4 digit');
//         });
//         expect(establishedOnField).toHaveAttribute('aria-invalid', 'true');
//       });

//       it('date of corporation field should have valid date', async () => {
//         await act(async () => {
//           render(<Certificate />);
//         });

//         const establishedOnField = screen.getByLabelText(/Date of Incorporation \/ stablishment/i);
//         userEvent.type(establishedOnField, '2020-01-02');

//         const submitButton = screen.getByRole('button', { name: /save & next/i });
//         userEvent.click(submitButton);

//         await waitFor(() => {
//           const errorMessages = screen
//             .getAllByTestId('error-message')
//             .map((err) => err.textContent);

//           expect(errorMessages).not.toContain('Invalid date / year must be 4 digit');
//         });

//         expect(establishedOnField).not.toHaveAttribute('aria-invalid');
//       });
//     });
//   });

//   describe('<LegalPerson />', () => {
//     it('should validate required field', async () => {
//       const initialState = {
//         stepper: {
//           currentStep: 2,
//           steps: [
//             { key: 1, status: 'complete', data: {} },
//             { key: 2, status: 'progress' }
//           ],
//           lastStep: null,
//           hasReachSubmitStep: false
//         }
//       };

//       await act(async () => {
//         render(<Certificate />, { preloadedState: initialState });
//       });

//       const submitButton = screen.getByRole('button', { name: /save & next/i });
//       userEvent.click(submitButton);

//       await waitFor(() => {
//         const errorMessages = screen.getAllByRole('alert').map((err) => err.textContent);
//         expect(errorMessages.length).toBe(5);
//       });
//     });
//   });

//   describe('<Contacts />', () => {
//     beforeAll(() => {
//       act(() => {
//         dynamicActivate('en');
//       });
//     });

//     // beforeEach(() => {
//     //   localStorage.setItem('certificateForm', JSON.stringify(certificateInitialValue));
//     // });

//     it('should validate required field', async () => {
//       const initialState = {
//         stepper: {
//           currentStep: 3,
//           steps: [
//             { key: 1, status: 'complete', data: {} },
//             { key: 2, status: 'progress' }
//           ],
//           lastStep: null,
//           hasReachSubmitStep: false
//         }
//       };
//       const contactFormValidationMessages = [
//         'Preferred name for email communication.',
//         'Email is required',
//         'Preferred name for email communication.',
//         'Email is required'
//       ];

//       await act(async () => {
//         render(<Certificate />, { preloadedState: initialState });
//       });

//       const submitButton = screen.getByRole('button', { name: /save & next/i });

//       userEvent.click(submitButton);

//       await waitFor(() => {
//         const errorMessages = screen.getAllByRole('alert').map((err) => err.textContent);
//         expect(errorMessages).toEqual(contactFormValidationMessages);
//       });
//     });
//   });

//   describe('<TrisaImplementation />', () => {
//     // beforeEach(() => {
//     //   localStorage.setItem('certificateForm', JSON.stringify(certificateInitialValue));
//     // });

//     it('should validate required field', async () => {
//       const initialState = {
//         stepper: {
//           currentStep: 4,
//           steps: [
//             { key: 1, status: 'complete', data: {} },
//             { key: 2, status: 'progress' }
//           ],
//           lastStep: null,
//           hasReachSubmitStep: false
//         }
//       };
//       const trisaImplementationFormValidationMessages = [
//         'TRISA endpoint is not valid',
//         'Common name should not contain special characters, no spaces and must have a dot(.) in it',
//         'TestNet and MainNet endpoints should not be the same',
//         'Common name should not contain special characters, no spaces and must have a dot(.) in it'
//       ];

//       await act(async () => {
//         render(<Certificate />, { preloadedState: initialState });
//       });

//       const submitButton = screen.getByRole('button', { name: /save & next/i });

//       userEvent.click(submitButton);

//       // await waitFor(() => {
//       //   const errorMessages = screen.getAllByRole('alert').map((err) => err.textContent);
//       //   expect(errorMessages).toEqual(trisaImplementationFormValidationMessages);
//       // });
//     });
//   });

//   describe('<TrixoQuestionnaire />', () => {
//     beforeEach(() => {
//       localStorage.setItem('certificateForm', JSON.stringify(certificateInitialValue));
//     });

//     it('should validate required field', async () => {
//       const initialState = {
//         stepper: {
//           currentStep: 4,
//           steps: [
//             { key: 1, status: 'complete', data: {} },
//             { key: 2, status: 'progress' }
//           ],
//           lastStep: null,
//           hasReachSubmitStep: false
//         }
//       };
//       const trixoQuestionnaireFormValidationMessages = [
//         'TRISA endpoint is not valid',
//         'Common name should not contain special characters, no spaces and must have a dot(.) in it',
//         'TestNet and MainNet endpoints should not be the same',
//         'Common name should not contain special characters, no spaces and must have a dot(.) in it'
//       ];

//       await act(async () => {
//         render(<Certificate />, { preloadedState: initialState });
//       });

//       const submitButton = screen.getByRole('button', { name: /save & next/i });

//       userEvent.click(submitButton);

//       // await waitFor(() => {
//       //   const errorMessages = screen.getAllByRole('alert').map((err) => err.textContent);
//       //   expect(errorMessages).toEqual(trixoQuestionnaireFormValidationMessages);
//       // });
//     });
//   });
// });

export {};
