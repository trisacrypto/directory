import { Meta, Story } from "@storybook/react";
import StatCard from ".";

interface StatCardProps {
  title: string;
  number: number;
}
export default {
  title: "components/StatCard",
  component: StatCard,
} as Meta;

const Template: Story = (args: StatCardProps) => <StatCard {...args} />;

export const Default = Template.bind({});
Default.args = {};
