import { Meta, Story } from "@storybook/react";
import LegalContact from ".";

export default {
  title: "components/Legal Contact",
  component: LegalContact,
} as Meta;

const Template: Story = (args) => <LegalContact {...args} />;

export const Default = Template.bind({});
Default.args = {};
