import { Meta, Story } from "@storybook/react";
import { InputProps } from "chakra-react-select";
import PhoneNumberInput from "./PhoneNumberInput";
import { Props } from "react-phone-number-input";

interface _Props extends Props<InputProps> {
  formHelperText?: string;
  controlId: string;
}

export default {
  title: "components/PhoneNumberInput",
  component: PhoneNumberInput,
  argTypes: { onChange: { action: "changed" } },
} as Meta<_Props>;

const Template: Story<any> = (args) => <PhoneNumberInput {...args} />;

export const Default = Template.bind({});
Default.args = {};
