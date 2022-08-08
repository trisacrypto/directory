import { Meta, Story } from '@storybook/react';
import StatCard, { StatCardProps } from '.';

export default {
  title: 'components/StatCard',
  component: StatCard
} as Meta;

const Template: Story<StatCardProps> = (args) => <StatCard {...args} />;

export const Default = Template.bind({});
Default.args = {};
