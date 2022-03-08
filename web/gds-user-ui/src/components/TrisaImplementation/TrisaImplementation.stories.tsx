import { Meta, Story } from "@storybook/react";
import TrisaImplementation from ".";

export default {
  title: "components/Trisa Implementation",
  component: TrisaImplementation,
} as Meta;

const Template: Story = (args) => <TrisaImplementation {...args} />;

export const Default = Template.bind({});
Default.args = {};
