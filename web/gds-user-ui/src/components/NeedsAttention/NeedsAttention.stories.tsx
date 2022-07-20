import { Meta, Story } from "@storybook/react";
import NeedsAttention, { NeedsAttentionProps } from ".";

export default {
  title: "components/NeedsAttention",
  component: NeedsAttention,
} as Meta;

const Template: Story<NeedsAttentionProps> = (args) => <NeedsAttention {...args} />;

export const Default = Template.bind({});
Default.args = {};
