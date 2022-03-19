import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import Addresses from '.';

type AddressesProps = {};

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
  title: 'components/Addresses',
  component: Addresses,
  decorators: [withRHF(false, defaultValues)]
} as Meta<AddressesProps>;

const Template: Story<AddressesProps> = (args) => <Addresses {...args} />;

export const Default = Template.bind({});
