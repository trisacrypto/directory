import { Meta, Story } from "@storybook/react";
import LegalPersonForm from "./LegalPersonForm";

type LegalPersonFormProps = {};

export default {
  title: "components/LegalPersonForm",
  component: LegalPersonForm,
} as Meta<LegalPersonFormProps>;

const Template: Story<LegalPersonFormProps> = (args) => (
  <LegalPersonForm {...args} />
);

export const Default = Template.bind({});
Default.args = {};
