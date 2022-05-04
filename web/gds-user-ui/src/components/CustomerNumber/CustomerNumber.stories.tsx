import { Meta, Story } from "@storybook/react";
import CustomerNumber from ".";

type CustomerNumberProps = {};

export default {
  title: "components/CustomerNumber",
  component: CustomerNumber,
} as Meta;

const Template: Story<CustomerNumberProps> = (args) => (
  <CustomerNumber {...args} />
);

export const Default = Template.bind({});
Default.args = {};
