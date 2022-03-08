import { Meta, Story } from "@storybook/react";
import AddressForm from ".";

type AddressFormProps = {};

export default {
  title: "components/AddressForm",
  component: AddressForm,
} as Meta<AddressFormProps>;

const Template: Story<AddressFormProps> = (args) => <AddressForm {...args} />;

export const Default = Template.bind({});
Default.args = {};
