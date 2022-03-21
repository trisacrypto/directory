import { Meta, Story } from "@storybook/react";
import TravelRuleProviders from ".";

export default {
  title: "components/TravelRuleProviders",
  component: TravelRuleProviders,
} as Meta;

const Template: Story = (args) => <TravelRuleProviders {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
