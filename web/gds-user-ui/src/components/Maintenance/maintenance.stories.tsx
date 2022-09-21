import { Meta, Story } from '@storybook/react';
import Maintenance from './';

export default {
  title: 'components/Maintenance',
  component: Maintenance
} as Meta;

const Template: Story = (args) => <Maintenance {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
