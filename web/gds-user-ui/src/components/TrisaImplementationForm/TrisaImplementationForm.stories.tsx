import { Meta, Story } from "@storybook/react";
import TrisaImplementationForm from ".";

export default {
  title: "components/Trisa Implementation Form",
  component: TrisaImplementationForm,
} as Meta;

const Template: Story = (args) => <TrisaImplementationForm {...args} />;

export const Default = Template.bind({});
Default.args = {};
