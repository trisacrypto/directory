import { Meta, Story } from '@storybook/react';
import { Control, UseFormRegister } from 'react-hook-form';
import AddressForm from '.';

type AddressFormProps = {
  control: Control;
  register: UseFormRegister<any>;
  name: string;
  rowIndex: number;
};

export default {
  title: 'components/AddressForm',
  component: AddressForm
} as Meta<AddressFormProps>;

const Template: Story<AddressFormProps> = (args) => <AddressForm {...args} />;

export const Default = Template.bind({});
Default.args = {};
