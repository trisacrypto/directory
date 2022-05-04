import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import NameIdentifiers from '.';

type NameIdentifiersProps = {};

const defaultValues = {
  name: {
    name_identifiers: [
      {
        legal_person_name: '',
        legal_person_name_identifier_type: '0'
      }
    ],
    local_name_identifiers: [
      {
        legal_person_name: '',
        legal_person_name_identifier_type: '0'
      }
    ],
    phonetic_name_identifiers: [
      {
        legal_person_name: '',
        legal_person_name_identifier_type: '0'
      }
    ]
  },
  geographic_addresses: [
    {
      address_type: 2,
      address_line: ['', '', ''],
      country: ''
    }
  ],
  customer_number: '',
  national_identification: {
    national_identifier: '',
    national_identifier_type: 0,
    country_of_issue: '',
    registration_authority: ''
  },
  country_of_registration: ''
};

export default {
  title: 'components/NameIdentifiers',
  component: NameIdentifiers
} as Meta<NameIdentifiersProps>;

const Template: Story<NameIdentifiersProps> = (args) => <NameIdentifiers {...args} />;

export const Default = Template.bind({});
Default.decorators = [withRHF(false, defaultValues)];
Default.args = {};
