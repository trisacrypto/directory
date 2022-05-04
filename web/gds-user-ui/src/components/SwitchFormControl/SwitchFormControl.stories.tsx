import { FormControlProps, SwitchProps } from "@chakra-ui/react";
import { Meta, Story } from "@storybook/react";
import SwitchFormControl from ".";

interface _FormControlProps extends FormControlProps {
  formHelperText?: string;
  controlId: string;
  label?: string;
  inputProps?: SwitchProps;
  name?: string;
  error?: string;
  direction?: "column" | "row";
}

export default {
  title: "components/ Switch Form Control",
  component: SwitchFormControl,
} as Meta<_FormControlProps>;

const Template: Story<_FormControlProps> = (args) => (
  <SwitchFormControl {...args} />
);

export const Standard = Template.bind({});
Standard.args = {
  label: "Enable email alerts?",
  formHelperText: "Form helper text",
};
