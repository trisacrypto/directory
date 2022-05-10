import { Meta, Story } from '@storybook/react';
import { SimpleDashboardLayout } from './SimpleDashboardLayout';

type SimpleDashboardLayout = {
  children: React.ReactNode;
};

export default {
  title: 'Layout/SimpleDashboardLayout',
  component: SimpleDashboardLayout
} as Meta;

const Template: Story<SimpleDashboardLayout> = (args) => <SimpleDashboardLayout {...args} />;

export const Standard = Template.bind({});
Standard.args = {
  children: 'this a children node'
};
