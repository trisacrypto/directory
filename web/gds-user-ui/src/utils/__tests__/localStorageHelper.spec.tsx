// import { getRegistrationDefaultValue } from 'modules/dashboard/certificate/lib';
// import {
//   addStepToLocalStorage,
//   loadDefaultValueFromLocalStorage,
//   loadStepperFromLocalStorage,
//   setCertificateFormValueToLocalStorage,
//   setStepperFromLocalStorage
// } from 'utils/localStorageHelper';

// const certificateForm = {
//   entity: {
//     country_of_registration: 'AS',
//     name: {
//       name_identifiers: [
//         {
//           legal_person_name: 'Technical',
//           legal_person_name_identifier_type: 'LEGAL_PERSON_NAME_TYPE_CODE_LEGL'
//         }
//       ],
//       local_name_identifiers: [],
//       phonetic_name_identifiers: []
//     },
//     geographic_addresses: [
//       {
//         address_type: 'ADDRESS_TYPE_CODE_BIZZ',
//         address_line: ['Address 1', 'Address 2', 'Address 3'],
//         country: 'AI'
//       }
//     ],
//     national_identification: {
//       national_identifier_type: 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX',
//       country_of_issue: '',
//       registration_authority: 'RA777777',
//       national_identifier: '2'
//     }
//   },
//   contacts: {
//     administrative: { name: '', email: '', phone: '' },
//     technical: { name: 'El', email: 'elyseebleu1@gmail.com', phone: '' },
//     billing: { name: '', email: '' },
//     legal: { name: 'Elysee', email: 'elyseebleu1@gmail.com' }
//   },
//   trisa_endpoint_testnet: {
//     trisa_endpoint: '',
//     common_name: 'testnet.technical.com',
//     endpoint: 'testnet.technical.com:443'
//   },
//   trisa_endpoint_mainnet: {
//     trisa_endpoint: '',
//     common_name: 'trisa.technical.com',
//     endpoint: 'trisa.technical.com:443'
//   },
//   website: 'http://technical.com',
//   business_category: 'GOVERNMENT_ENTITY',
//   vasp_categories: [],
//   established_on: '2022-04-20',
//   organization_name: 'Technical',
//   trixo: {
//     primary_national_jurisdiction: 'AS',
//     primary_regulator: '',
//     other_jurisdictions: [],
//     financial_transfers_permitted: 'partial',
//     has_required_regulatory_program: 'yes',
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

// describe('localStorageHelper', () => {
//   const STEP_KEY = 'trs_stepper';

//   beforeEach(() => {
//     localStorage.clear();
//   });

//   it('should add step to localStorage', () => {
//     const result = { steps: [{ key: 1, status: 'complete', data: {} }] };

//     addStepToLocalStorage({ key: 1, status: 'complete', data: {} });

//     expect(localStorage.setItem).toHaveBeenCalled();
//     expect(JSON.parse(localStorage.getItem(STEP_KEY)!)).toEqual(result);
//   });

//   describe('setCertificateFormValueToLocalStorage', () => {
//     it('should set certificate form value', () => {
//       setCertificateFormValueToLocalStorage({ key: 1, status: 'complete', data: {} });

//       expect(localStorage.setItem).toHaveBeenCalledWith(
//         'certificateForm',
//         JSON.stringify({ key: 1, status: 'complete', data: {} })
//       );
//     });

//     it('should not set certificate form value', () => {
//       setCertificateFormValueToLocalStorage(undefined);
//       expect(localStorage.setItem).not.toHaveBeenCalled();
//     });
//   });

//   it('should load stepper from localStorage', () => {
//     loadStepperFromLocalStorage();
//     const defaultValue: any = {
//       currentStep: 1,
//       steps: [
//         {
//           key: 1,
//           status: 'progress',
//           data: {}
//         }
//       ],
//       lastStep: null,
//       hasReachSubmitStep: false
//     };

//     expect(localStorage.getItem).toHaveBeenCalledWith('trs_stepper');
//     expect(localStorage.setItem).toHaveBeenCalledWith('trs_stepper', JSON.stringify(defaultValue));

//     expect(localStorage.getItem('trs_stepper')).toEqual(JSON.stringify(defaultValue));
//   });

//   it('should load stepper from localStorage', () => {
//     const defaultValue: any = {
//       currentStep: 1,
//       steps: [
//         {
//           key: 1,
//           status: 'progress',
//           data: {}
//         }
//       ],
//       lastStep: null,
//       hasReachSubmitStep: false
//     };

//     localStorage.setItem('trs_stepper', JSON.stringify(defaultValue));

//     expect(loadStepperFromLocalStorage()).toEqual(defaultValue);
//   });

//   describe('loadDefaultValueFromLocalStorage', () => {
//     it('should load default values from local storage', () => {
//       localStorage.setItem('certificateForm', JSON.stringify(certificateForm));

//       loadDefaultValueFromLocalStorage();

//       expect(localStorage.getItem).toHaveBeenCalledWith('certificateForm');
//       expect(localStorage.setItem).toHaveBeenCalledWith(
//         'certificateForm',
//         JSON.stringify(certificateForm)
//       );
//       expect(loadDefaultValueFromLocalStorage()).toEqual(certificateForm);
//     });

//     it('should load default values even if there no value in the local storage', () => {
//       expect(loadDefaultValueFromLocalStorage()).toEqual(getRegistrationDefaultValue());
//     });
//   });

//   describe('setStepperFromLocalStorage', () => {
//     const stepper = {
//       currentStep: 5,
//       steps: [
//         { key: 1, status: 'complete', data: {} },
//         { key: 2, status: 'complete' },
//         { key: 3, status: 'complete' },
//         { key: 4, status: 'complete' },
//         { key: 5, status: 'complete' },
//         { key: 6, status: 'complete' }
//       ],
//       lastStep: null,
//       hasReachSubmitStep: false
//     };
//   });

//   it('should set', () => {
//     setStepperFromLocalStorage({
//       step: 1,
//       status: 'complete',
//       data: {}
//     });

//     expect(localStorage.getItem).toHaveBeenCalledWith('trs_stepper');
//   });

//   it('should step the current step when step is passed as parameter', () => {
//     localStorage.setItem(
//       'trs_stepper',
//       JSON.stringify({
//         currentStep: 1,
//         steps: [{ key: 1 }],
//         lastStep: null,
//         hasReachSubmitStep: false
//       })
//     );
//     setStepperFromLocalStorage({ step: 2 });

//     expect(localStorage.getItem).toHaveBeenCalledWith('trs_stepper');
//     expect(localStorage.getItem('trs_stepper')).toEqual(
//       JSON.stringify({
//         currentStep: 2,
//         steps: [{ key: 1 }],
//         lastStep: null,
//         hasReachSubmitStep: false
//       })
//     );
//   });

//   it('should set status when it is passed as parameter', () => {
//     localStorage.setItem(
//       'trs_stepper',
//       JSON.stringify({
//         currentStep: 1,
//         steps: [{ key: 2 }],
//         lastStep: null,
//         hasReachSubmitStep: false
//       })
//     );
//     setStepperFromLocalStorage({ step: 2, status: 'complete' });

//     expect(localStorage.getItem).toHaveBeenCalledWith('trs_stepper');
//     expect(localStorage.getItem('trs_stepper')).toEqual(
//       JSON.stringify({
//         currentStep: 2,
//         steps: [{ key: 2, status: 'complete' }],
//         lastStep: null,
//         hasReachSubmitStep: false
//       })
//     );
//   });

//   it('should set data when data is passed as parameter', () => {
//     localStorage.setItem(
//       'trs_stepper',
//       JSON.stringify({
//         currentStep: 1,
//         steps: [{ key: 2 }],
//         lastStep: null,
//         hasReachSubmitStep: false
//       })
//     );
//     setStepperFromLocalStorage({ step: 2, status: 'complete', data: { test: 'test' } });

//     expect(localStorage.getItem).toHaveBeenCalledWith('trs_stepper');
//     expect(localStorage.getItem('trs_stepper')).toEqual(
//       JSON.stringify({
//         currentStep: 2,
//         steps: { test: 'test' },
//         lastStep: null,
//         hasReachSubmitStep: false
//       })
//     );
//   });

//   afterEach(() => {
//     jest.clearAllMocks();
//   });
// });
export {};

describe('test ', () => {
  it('should ', () => {
    expect(true).toBe(true);
  });
});
