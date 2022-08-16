import { Meta, Story } from '@storybook/react';
import StatisticCard, { StatisticCardProps } from './StatisticCard';

export default {
  title: 'components/StatisticCard',
  component: StatisticCard
} as Meta<StatisticCardProps>;

const Template: Story<StatisticCardProps> = (args) => <StatisticCard {...args} />;

export const Default = Template.bind({});
Default.args = {
  title: 'Network Status',
  total: 0
};
