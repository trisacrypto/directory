import { ButtonProps } from "@chakra-ui/react";
import { Meta, Story } from "@storybook/react";
import FormButton from "./FormButton";

interface _ButtonProps extends ButtonProps {}

export default {
  title: "components/FormButton",
  component: FormButton,
} as Meta<_ButtonProps>;

const Template: Story<_ButtonProps> = (args) => <FormButton {...args} />;

export const Default = Template.bind({});
Default.args = {
  children: "Button",
};
