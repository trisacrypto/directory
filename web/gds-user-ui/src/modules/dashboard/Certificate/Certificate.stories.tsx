import { Meta, Story } from '@storybook/react';
import Certificate from './Certificate';

export default {
  title: 'modules/Certificate',
  component: Certificate
} as Meta;

const Template: Story = (args) => <Certificate {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
