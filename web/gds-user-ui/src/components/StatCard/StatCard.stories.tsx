import { Meta, Story } from "@storybook/react";
import StatCard from ".";

export default {
  title: "components/StatCard",
  component: StatCard,
} as Meta;

const Template: Story = (args) => <StatCard {...args} />;

export const Default = Template.bind({});
Default.args = {};
