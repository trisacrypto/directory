import { Meta, Story } from "@storybook/react";
import TechnicalResources from ".";

export default {
  title: "components/TechnicalResources",
  component: TechnicalResources,
} as Meta;

const Template: Story = (args) => <TechnicalResources {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
