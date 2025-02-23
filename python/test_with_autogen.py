import asyncio
from typing import List
from autogen_agentchat.agents import AssistantAgent
from autogen_agentchat.ui import Console
from autogen_agentchat.conditions import TextMentionTermination, MaxMessageTermination
from autogen_agentchat.teams import RoundRobinGroupChat, SelectorGroupChat
from autogen_ext.models.openai import (
    OpenAIChatCompletionClient,
    AzureOpenAIChatCompletionClient,
)
from autogen_ext.tools.mcp import mcp_server_tools, StdioMcpToolAdapter, StdioServerParams
from dotenv import load_dotenv
import os
import sys
from pathlib import Path
import httpx
from rich.console import Console as RichConsole

def print_mcp_tools(tools: List[StdioMcpToolAdapter]) -> None:
    """Print available MCP tools and their parameters in a formatted way."""
    console = RichConsole()
    console.print("\n[bold blue]📦 Loaded MCP Tools:[/bold blue]\n")

    for tool in tools:
        # Tool name and description
        console.print(f"[bold green]🔧 {tool.schema.get('name', 'Unnamed Tool')}[/bold green]")
        if description := tool.schema.get('description'):
            console.print(f"[italic]{description}[/italic]\n")

        # Parameters section
        if params := tool.schema.get('parameters'):
            console.print("[yellow]Parameters:[/yellow]")
            if properties := params.get('properties', {}):
                required_params = params.get('required', [])
                for prop_name, prop_details in properties.items():
                    required_mark = "[red]*[/red]" if prop_name in required_params else ""
                    param_type = prop_details.get('type', 'any')
                    console.print(f"  • [cyan]{prop_name}{required_mark}[/cyan]: {param_type}")
                    if param_desc := prop_details.get('description'):
                        console.print(f"    [dim]{param_desc}[/dim]")

        console.print("─" * 60 + "\n")



async def main() -> None:
   
    argo_system_mcp_server = StdioServerParams(
        command="/workspaces/mcp-argo-server/bin/mcp-argo-server",        
    )

    file_system_tools = await mcp_server_tools(argo_system_mcp_server)
    print_mcp_tools(file_system_tools)


asyncio.run(main())
