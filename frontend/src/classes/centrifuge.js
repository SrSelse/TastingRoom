import { Centrifuge, UnauthorizedError } from 'centrifuge';

export class OPCentrifuge {
  centrifuge = null;
  refCounts = new Map();

  async init(token = null) {
    if (this.centrifuge) {
      // await this.centrifuge.ready(2500);
      return;
    }

    let protocol = 'ws';
    let url = import.meta.env.VITE_CENTRIFUGO_WS;

    if (window.location.protocol === 'https:') {
      protocol = 'wss';
      url = import.meta.env.VITE_CENTRIFUGO_WSS;
    }

    this.centrifuge = new Centrifuge(`${protocol}://${url}/connection/websocket`, {
      getToken: this.generateConnectionToken
    });

    if (token) {
      this.centrifuge.setToken(token);
    }

    this.centrifuge.connect();
  }

  async generateConnectionToken() {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/broadcasting/connect`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      if (!res.ok) {
        if (res.status === 403) {
          // Return special error to not proceed with token refreshes, client will be disconnected.
          throw new UnauthorizedError();
        }
        // Any other error thrown will result into token refresh re-attempts.
        throw new Error(`Unexpected status code ${res.status}`);
      }
      const data = await res.json();
      return data.token;
    } catch (e) {
      console.log(e);
    }
  }

  async getSubscriptionToken(ctx) {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/broadcasting/auth`, {
        method: "POST",
        body: JSON.stringify(ctx),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await res.json();
      return data.token;
    } catch (e) {
      if (e.response && e.response.status === 403) {
        console.error(e)
        return "";
      }

      throw e;
    }
  }

  newSubscription (name, options) {
    // Reuse existing subscription — multiple components may listen to the
    // same channel (e.g. AdminMode + SatelliteMode on the same beer).
    let sub = this.centrifuge.getSubscription(name);

    if (!sub) {
      const subscriptionOptions = {...options};

      // Inget behov av subscription token för user channels, men bör kunna hanteras snyggare?
      if (name.indexOf('#') === -1) {
        subscriptionOptions.getToken = (ctx) => this.getSubscriptionToken(ctx);
      }

      sub = this.centrifuge.newSubscription(name, {
        ...subscriptionOptions,
      });
    }

    this.refCounts.set(name, (this.refCounts.get(name) || 0) + 1);
    return sub;
  }

  removeSubscription(sub) {
    if (!sub) {
      return;
    }

    // Only tear down the channel when the last listener leaves.
    const count = (this.refCounts.get(sub.channel) || 1) - 1;
    if (count > 0) {
      this.refCounts.set(sub.channel, count);
      return;
    }

    this.refCounts.delete(sub.channel);
    sub.unsubscribe();
    this.centrifuge.removeSubscription(sub);
  }
}

export const centrifuge = new OPCentrifuge();

