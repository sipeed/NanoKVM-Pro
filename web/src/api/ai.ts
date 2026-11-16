import { http } from '@/lib/http.ts';

export type AIConfig = {
  mcpEnabled: boolean;
  mcpApiKeyMasked?: string;
};

export const getAIConfig = () => http.get('/api/ai/config');
export const setAIConfig = (mcpEnabled: boolean) => http.post('/api/ai/config', { mcpEnabled });
export const getMCPAPIKey = () => http.post('/api/ai/mcp/key');
export const copyMCPAPIKey = () => http.post('/api/ai/mcp/key/copy');
export const regenerateMCPAPIKey = () => http.post('/api/ai/mcp/key/regenerate');
