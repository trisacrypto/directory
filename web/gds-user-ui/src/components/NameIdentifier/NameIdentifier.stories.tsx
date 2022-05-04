import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import { Control, FieldValues, RegisterOptions } from 'react-hook-form';
import NameIdentifier from '.';

type NameIdentifierProps = {
  name: string;
  description: string;
  controlId: string;
  register: RegisterOptions;
  control: Control<FieldValues, any>;
  heading: string;
};

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
  title: 'components/NameIdentifer',
  component: NameIdentifier
} as Meta<NameIdentifierProps>;

const Template: Story<NameIdentifierProps> = (args) => (
  <NameIdentifier {...args} name="name.name_identifiers" />
);

export const Default = Template.bind({});
Default.decorators = [withRHF(false, defaultValues)];
Default.args = {
  description: 'Enter at least one geographic address.',
  heading: 'Name Identifiers'
};
