import { Meta, Story } from "@storybook/react";
import TrixoQuestionnaireForm from ".";

export default {
  title: "components/TrixoQuestionnaireForm",
  component: TrixoQuestionnaireForm,
} as Meta;

const Template: Story = (args) => <TrixoQuestionnaireForm {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
