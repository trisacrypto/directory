import { Meta, Story } from '@storybook/react';
import { GroupBase, OptionsOrGroups, Props } from 'chakra-react-select';
import SelectFormControl from '.';

const options = [
  { value: 'AL', label: 'Alabama' },
  { value: 'AK', label: 'Alaska' },
  { value: 'AS', label: 'American Samoa' }
];

interface _FormControlProps extends Props {
  formHelperText?: string;
  controlId: string;
  label?: string;
  name?: string;
  placeholder?: string;
  options?: OptionsOrGroups<unknown, GroupBase<unknown>>;
}

export default {
  title: 'components/SelectFormControl',
  component: SelectFormControl
} as Meta<_FormControlProps>;

const Template: Story<_FormControlProps> = (args) => <SelectFormControl {...args} />;

export const Default = Template.bind({});
Default.args = {
  options,
  formHelperText: 'Choose one country',
  label: 'Country',
  isMulti: false
};

export const Invalid = Template.bind({});
Invalid.args = {
  ...Default.args,
  formHelperText: 'Pick at least one country',
  isInvalid: true
};
