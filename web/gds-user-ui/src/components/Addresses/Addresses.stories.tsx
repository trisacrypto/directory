import { Meta, Story } from "@storybook/react";
import Addresses from "./";

type AddressesProps = {};

export default {
  title: "components/Addresses",
  component: Addresses,
} as Meta<AddressesProps>;

const Template: Story<AddressesProps> = (args) => <Addresses {...args} />;

export const Default = Template.bind({});
