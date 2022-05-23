/* eslint-disable react/display-name */
import { Trans } from '@lingui/react';
import { action } from '@storybook/addon-actions';
import { StoryFnReactReturnType } from '@storybook/react/dist/ts3.9/client/preview/types';
import { VFC, ReactNode, FC } from 'react';
import { FormProvider, useForm } from 'react-hook-form';

const StorybookFormProvider: VFC<{
  children: ReactNode;
  defaultValues?: {
    [x: string]: any;
  };
}> = ({ children, defaultValues }) => {
  const methods = useForm({
    defaultValues: defaultValues || {
      entity: {
        name: {
          name_identifiers: [],
          local_name_identifiers: [],
          phonetic_name_identifiers: [
            {
              legal_person_name: '',
              legal_person_name_identifier_type: 0
            }
          ]
        },
        geographic_addresses: [
          {
            address_type: 2,
            address_line: ['', '', ''],
            country: ''
          },
          {
            address_type: 2,
            address_line: ['', '', ''],
            country: ''
          }
        ],
        customer_number: '',
        national_identification: {
          national_identifier: '',
          national_identifier_type: 9,
          country_of_issue: '',
          registration_authority: ''
        },
        country_of_registration: ''
      },
      contacts: {
        technical: {
          name: 'Koffi',
          email: 'koffi@gmail.com',
          phone: '21803u485'
        },
        legal: {
          name: '',
          email: '',
          phone: ''
        },
        administrative: {
          name: '',
          email: '',
          phone: ''
        },
        billing: {
          name: 'Elysee ',
          email: 'elyseebleu2gmail.com',
          phone: ''
        }
      },
      trisa_endpoint: '',
      common_name: '',
      website: '',
      business_category: 4,
      vasp_categories: ['DEX', 'P2P', 'Kiosk'],
      established_on: '11111-11-11',
      trixo: {
        primary_national_jurisdiction: '',
        primary_regulator: '',
        other_jurisdictions: [
          {
            country: '',
            regulator_name: ''
          }
        ],
        financial_transfers_permitted: 'partial',
        has_required_regulatory_program: 'partial',
        conducts_customer_kyc: true,
        kyc_threshold: 0,
        kyc_threshold_currency: 'XCD',
        must_comply_travel_rule: false,
        applicable_regulations: ['FATF Recommendation 16'],
        compliance_threshold: 0,
        compliance_threshold_currency: 'VES',
        must_safeguard_pii: false,
        safeguards_pii: false
      }
    }
  });

  return (
    <FormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(action('[React Hooks Form] Submit'))}>{children}</form>
    </FormProvider>
  );
};

StorybookFormProvider.displayName = 'StorybookFormProvider';

export const withRHF =
  (showSubmitButton: boolean, defaultValues?: Record<string, any>) =>
  (Story: FC): StoryFnReactReturnType =>
    (
      <StorybookFormProvider defaultValues={defaultValues}>
        <Story />
        {showSubmitButton && (
          <button type="submit">
            <Trans id="Submit">Submit</Trans>{' '}
          </button>
        )}
      </StorybookFormProvider>
    );
