import { Meta, Story } from "@storybook/react";
import NameIdentifier from ".";

type NameIdentifierProps = {
  name: string;
  description: string;
};

export default {
  title: "components/NameIdentifer",
  component: NameIdentifier,
} as Meta<NameIdentifierProps>;

const Template: Story<NameIdentifierProps> = (args) => (
  <NameIdentifier {...args} />
);

export const Default = Template.bind({});
Default.args = {
  name: "Name Identifiers",
  description: "Enter at least one geographic address.",
};
