import { Meta, Story } from "@storybook/react";
import LegalPerson from "./LegalPerson";

export default {
  title: "components/Legal Person",
  component: LegalPerson,
} as Meta;

const Template: Story = (args) => <LegalPerson {...args} />;

export const Default = Template.bind({});
Default.args = {};
