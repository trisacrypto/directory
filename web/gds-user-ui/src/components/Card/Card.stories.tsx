import { Meta, Story } from "@storybook/react";
import Card from ".";

export default {
  title: "components/Card",
  component: Card,
} as Meta;

const Template: Story = (args) => <Card {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
