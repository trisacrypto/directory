import { Meta, Story } from "@storybook/react";
import YourImplementation from ".";

export default {
  title: "components/YourImplementation",
  component: YourImplementation,
} as Meta;

const Template: Story = (args) => <YourImplementation {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
