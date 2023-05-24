import { Story } from '@storybook/react';

import BasicDetailsForm from '.';

type BasicDetailsFormProps = {
  data: any;
};

export default {
  title: 'components/BasicDetailsForm',
  component: BasicDetailsForm
};

const Template: Story<BasicDetailsFormProps> = (args) => <BasicDetailsForm {...args} />;

export const Default = Template.bind({});
Default.args = {};
