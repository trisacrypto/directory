import { Meta, Story } from '@storybook/react';
import ChooseAnOrganization, { ChooseAnAccountProps } from '.';

export default {
  title: 'components/ChooseAnAccount',
  component: ChooseAnOrganization
} as Meta<ChooseAnAccountProps>;

const Template: Story<ChooseAnAccountProps> = (args) => <ChooseAnOrganization {...args} />;

export const Standard = Template.bind({});
Template.args = {};
