import { Meta, Story } from "@storybook/react";
import NeedsAttention from ".";

export default {
  title: "components/NeedsAttention",
  component: NeedsAttention,
} as Meta;

const Template: Story = (args) => <NeedsAttention {...args} />;

export const Default = Template.bind({});
Default.args = {};
