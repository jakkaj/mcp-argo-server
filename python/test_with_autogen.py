import asyncio
from typing import List
from autogen_agentchat.agents import AssistantAgent
from autogen_agentchat.ui import Console
from autogen_agentchat.conditions import TextMentionTermination, MaxMessageTermination
from autogen_agentchat.teams import RoundRobinGroupChat, SelectorGroupChat
from autogen_core import CancellationToken
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
from pydantic import BaseModel
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

class ArgoSubmitConfigModel(BaseModel):
    manifest: str  # You can refine this type if needed
    namespace: str
    wait: bool = False

class ArgoStatusConfigModel(BaseModel):
    name: str
    namespace: str

async def main() -> None:
   
    argo_system_mcp_server = StdioServerParams(
        command="/workspaces/mcp-argo-server/bin/mcp-argo-server",        
    )

    argo_tools = await mcp_server_tools(argo_system_mcp_server)
    print_mcp_tools(argo_tools)

    token = CancellationToken()

    argo_manifest = ""
    # load manifest string from ../kube/argo-hello-world.yaml
    with open("kube/argo-hello-world.yaml", "r") as f:
        argo_manifest = f.read()

    if not argo_manifest:
        print("argo_manifest is empty")
        sys.exit(1)

    print(argo_tools)
     #### Wait for the job to finish and get the output
    argo_config_wait = ArgoSubmitConfigModel(manifest=argo_manifest, namespace="argo", wait=True)
    res4 = await argo_tools[0].run(argo_config_wait, token)    
    print(res4)
    # waited_name = res4[0:1][0]
    # argo_status_config = ArgoStatusConfigModel(name=waited_name.text, namespace="argo")
    # waitedResult = await argo_tools[1].run(argo_status_config, token)
    # print(waitedResult)


    ### Don't wait auto, do it this side
    argo_config = ArgoSubmitConfigModel(manifest=argo_manifest, namespace="argo")

    res = await argo_tools[0].run(argo_config, token)
    
    print(res)

    name = res[-1:][0]

    print (name.text)

    argo_status_config = ArgoStatusConfigModel(name=name.text, namespace="argo")
    
    while True:
        res2 = await argo_tools[2].run(argo_status_config, token)
        phase = res2[-1:][0].text
        print(phase)
        if phase == "Succeeded":
            break
        await asyncio.sleep(2)

    print(res2)

    res3 = await argo_tools[1].run(argo_status_config, token)
    print(res3)

asyncio.run(main())
