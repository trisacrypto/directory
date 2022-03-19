import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import LegalPerson from '.';

const defaultValues = {
  name: {
    name_identifiers: [
      {
        legal_person_name: '',
        legal_person_name_identifier_type: '0'
      }
    ],
    local_name_identifiers: [],
    phonetic_name_identifiers: []
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
  title: 'components/Legal Person',
  component: LegalPerson,
  decorators: [withRHF(false, defaultValues)]
} as Meta;

const Template: Story = (args) => <LegalPerson {...args} />;

export const Default = Template.bind({});
Default.args = {};
