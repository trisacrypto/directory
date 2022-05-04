import { InputProps } from "@chakra-ui/react";
import { Story } from "@storybook/react";
import InputFormControl from ".";

interface _FormControlProps {
  formHelperText?: string;
  controlId: string;
  label?: string;
  inputProps?: InputProps;
  isInvalid?: boolean;
}

export default {
  title: "components/InputFormControl",
  component: InputFormControl,
};

const Template: Story<_FormControlProps> = (args) => {
  return <InputFormControl {...args} />;
};

export const Default = Template.bind({});
Default.args = {
  label: "Email",
  formHelperText: "Your form helper text",
  inputProps: {
    type: "email",
    placeholder: "your placeholder",
  },
};

export const Invalid = Template.bind({});
Invalid.args = {
  ...Default.args,
  formHelperText: "Email is required",
  isInvalid: true,
};
